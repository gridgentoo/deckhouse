/*
Copyright 2022 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hooks

import (
	"context"
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	d8config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/modules"
	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/go_lib/dependency/k8s"
)

/**
This hook switches deckhouse-controller to use a number of typed DeckhouseConfig custom
resources and a managed ConfigMap/deckhouse-generated-config-do-not-edit
instead one untyped ConfigMap/deckhouse.
*/

// Use order:1 to run before all global hooks.
var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnStartup: &go_hook.OrderedConfig{Order: 1},
}, dependency.WithExternalDependencies(migrateOrSyncModuleConfigs))

// migrateOrSyncModuleConfigs runs on deckhouse-controller startup
// as early as possible and do two things:
// - switches deckhouse-controller from configuration via ConfigMap/deckhouse
//   that is managed by deckhouse and by user to configuration via
//   DeckhouseConfig objects that managed by user so can be stored in Git.
// - synchronize DeckhouseConfig objects content to intermediate
//   ConfigMap/deckhouse-generated-config-do-not-edit.
func migrateOrSyncModuleConfigs(input *go_hook.HookInput, dc dependency.Container) error {
	err := modules.Registry().Init(modules.DefaultGlobalHooksDir, modules.DefaultModulesDir)
	if err != nil {
		return fmt.Errorf("initialize modules registry: %v", err)
	}

	kubeClient, err := dc.GetK8sClient()
	if err != nil {
		return fmt.Errorf("cannot init Kubernetes client: %v", err)
	}

	// Detect if Deployment/deckhouse is used ConfigMap/deckhouse-generated-config-do-not-edit.
	isGeneratedMapInUse, err := isGeneratedMapInUse(kubeClient)
	if err != nil {
		return fmt.Errorf("get deploy/deckhouse: %v", err)
	}

	// Get ConfigMap/deckhouse-generated-config-do-not-edit.
	hasGeneratedCM := true
	generatedCM, err := d8config.GetGeneratedConfigMap(kubeClient)
	if err != nil {
		if !k8errors.IsNotFound(err) {
			return fmt.Errorf("get generated ConfigMap: %v", err)
		}
		// NotFound error is occurred.
		hasGeneratedCM = false
	}

	// List DeckhouseConfig objects.
	allConfigs, err := d8config.GetAllConfigs(kubeClient)
	if err != nil {
		return fmt.Errorf("get all settings: %v", err)
	}

	if !isGeneratedMapInUse {
		// Migrate ConfigMap/deckhouse to ModuleConfigs
		// - create DeckhouseConfig resources for each module
		// - create new generated ConfigMap
		// - update deploy/deckhouse to use new ConfigMap.
		input.LogEntry.Infof("Deployment/deckhouse is not migrated. Migrate to separate module configs.")
		return migrateToModuleConfigs(input, kubeClient)
	}

	// Create absent DeckhouseConfig objects from generated ConfigMap and
	// then sync generated ConfigMap from DeckhouseConfig objects.
	// Note: DeckhouseConfig objects are controlled from Git, so auto-deletion is not an option.
	// The compromise here is to use generated ConfigMap to re-create DeckhouseConfig objects for known modules only.
	input.LogEntry.Infof("Generated cm: %v, module configs: %d. Run sync.", hasGeneratedCM, len(allConfigs))
	return syncModuleConfigs(input, generatedCM, allConfigs)
}

func migrateToModuleConfigs(input *go_hook.HookInput, kubeClient k8s.Client) error {
	possibleNames := modules.Registry().GetPossibleNames()
	// Migrate cm/deckhouse to DeckhouseConfig resources and a generated ConfigMap.
	objs, err := d8config.MigrateConfigMapToModuleConfigs(kubeClient, possibleNames)
	if err != nil {
		return err
	}

	input.LogEntry.Infof("Create %d DeckhouseConfig objects", len(objs))
	for _, obj := range objs {
		input.PatchCollector.Create(obj, object_patch.UpdateIfExists())
	}

	cmData, err := d8config.SyncFromDeckhouseConfigs(objs)
	if err != nil {
		return err
	}
	cm := d8config.GeneratedConfigMap(cmData)
	input.LogEntry.Infof("Create Config/%s", cm.Name)
	input.PatchCollector.Create(cm, object_patch.UpdateIfExists())

	// Patch deploy/deckhouse to use generated ConfigMap for config values.
	switchDeckhouseToGeneratedConfigMap := func(u *unstructured.Unstructured) (*unstructured.Unstructured, error) {
		var depl appsv1.Deployment
		err := sdk.FromUnstructured(u, &depl)
		if err != nil {
			return nil, err
		}

		envs := depl.Spec.Template.Spec.Containers[0].Env
		newEnvs := make([]v1.EnvVar, 0)
		for _, envVar := range envs {
			if envVar.Name == "ADDON_OPERATOR_CONFIG_MAP" {
				newEnvs = append(newEnvs, v1.EnvVar{
					Name:  "ADDON_OPERATOR_CONFIG_MAP",
					Value: cm.Name,
				})
			} else {
				newEnvs = append(newEnvs, envVar)
			}
		}

		depl.Spec.Template.Spec.Containers[0].Env = newEnvs
		return sdk.ToUnstructured(&depl)
	}
	input.LogEntry.Infof("Update deploy/deckhouse to use Config/%s", cm.Name)
	input.PatchCollector.Filter(switchDeckhouseToGeneratedConfigMap, "apps/v1", "Deployment", d8config.DeckhouseNS, "deckhouse")
	// Deckhouse restarts here.
	return nil
}

func syncModuleConfigs(input *go_hook.HookInput, generatedCM *v1.ConfigMap, allConfigs []*d8config_v1.DeckhouseConfig) error {
	cmData, err := d8config.SyncFromDeckhouseConfigs(allConfigs)
	if err != nil {
		return err
	}
	cm := d8config.GeneratedConfigMap(cmData)
	input.PatchCollector.Create(cm, object_patch.UpdateIfExists())

	if generatedCM == nil || len(generatedCM.Data) == 0 {
		return nil
	}

	// Log deleted sections.
	for name := range generatedCM.Data {
		if _, has := cmData[name]; !has {
			input.LogEntry.Warnf("Seems DeckhouseConfig/%s was deleted. Delete section '%s' from cm/%s.", name, name, cm.Name)
		}
	}

	return nil
}

// isGeneratedMapInUse detects if Deployment/deckhouse is already
// migrated to use generated ConfigMap.
func isGeneratedMapInUse(klient k8s.Client) (bool, error) {
	depl, err := klient.AppsV1().Deployments(d8config.DeckhouseNS).Get(context.TODO(), "deckhouse", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if depl == nil {
		return false, fmt.Errorf("no object")
	}
	envs := depl.Spec.Template.Spec.Containers[0].Env
	for _, envVar := range envs {
		if envVar.Name == "ADDON_OPERATOR_CONFIG_MAP" {
			return envVar.Value == d8config.GeneratedConfigMapName, nil
		}
	}
	return false, fmt.Errorf("env ADDON_OPERATOR_CONFIG_MAP not found in deployment spec")
}

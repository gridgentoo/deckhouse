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
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"

	d8config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	"github.com/deckhouse/deckhouse/modules/019-deckhouse-config/hooks/internal"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/deckhouse-config/updater",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:                   "configs",
			ApiVersion:             "deckhouse.io/v1",
			Kind:                   "DeckhouseConfig",
			WaitForSynchronization: pointer.BoolPtr(true),
			FilterFunc:             filterDeckhouseConfigs,
		},
		{
			Name:       "generated-cm",
			ApiVersion: "v1",
			Kind:       "ConfigMap",
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{d8config.DeckhouseNS},
				},
			},
			NameSelector: &types.NameSelector{
				MatchNames: []string{d8config.GeneratedConfigMapName},
			},
			ExecuteHookOnEvents:          pointer.BoolPtr(false),
			ExecuteHookOnSynchronization: pointer.BoolPtr(false),
			FilterFunc:                   filterGeneratedConfigMap,
		},
	},
}, updateGeneratedConfigMap)

// filterModuleSettings returns spec for DeckhouseConfig objects.
func filterDeckhouseConfigs(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var cfg d8config_v1.DeckhouseConfig

	err := sdk.FromUnstructured(unstructured, &cfg)
	if err != nil {
		return nil, err
	}

	// Extract name and spec into empty DeckhouseConfig.
	return &d8config_v1.DeckhouseConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: cfg.Name,
		},
		Spec: cfg.Spec,
	}, nil
}

type configData map[string]string

// filterGeneratedConfigMap returns Data field for ConfigMap.
func filterGeneratedConfigMap(unstructured *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var cm v1.ConfigMap

	err := sdk.FromUnstructured(unstructured, &cm)
	if err != nil {
		return nil, err
	}

	return configData(cm.Data), nil
}

// updateGeneratedConfigMap converts specs from ModuleSettings resources into ConfigMap data.
// TODO add deletion threshold.
func updateGeneratedConfigMap(input *go_hook.HookInput) error {
	namesSet, err := internal.SetFromArrayValue(input, PossibleNamesPath)
	if err != nil {
		return err
	}

	allConfigs := internal.KnownConfigsFromSnapshot(input.Snapshots["configs"], namesSet)

	// TODO add logic to ignore specific module configs?
	// TODO add logic to postpone deletion sections in ConfigMap?

	for _, cfg := range allConfigs {
		err := d8config.ValidateValues(cfg)
		if err != nil {
			return err
		}
	}

	cmData, err := d8config.SyncFromDeckhouseConfigs(allConfigs)
	if err != nil {
		return fmt.Errorf("convert DeckhouseConfig objects to ConfigMap: %s", err)
	}

	cm := d8config.GeneratedConfigMap(cmData)
	input.PatchCollector.Create(cm, object_patch.UpdateIfExists())

	return nil
}

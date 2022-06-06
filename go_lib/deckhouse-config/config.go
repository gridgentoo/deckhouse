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

package deckhouse_config

import (
	"context"
	"fmt"

	kcm "github.com/flant/addon-operator/pkg/kube_config_manager"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	"github.com/deckhouse/deckhouse/go_lib/dependency/k8s"
)

// GetAllConfigs returns all DeckhouseConfig objects.
func GetAllConfigs(kubeClient k8s.Client) ([]*d8config_v1.DeckhouseConfig, error) {
	gvr := d8config_v1.GroupVersionResource()
	unstructuredObjs, err := kubeClient.Dynamic().Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	objs := make([]*d8config_v1.DeckhouseConfig, 0, len(unstructuredObjs.Items))
	for _, unstructured := range unstructuredObjs.Items {
		var obj d8config_v1.DeckhouseConfig

		err := sdk.FromUnstructured(&unstructured, &obj)
		if err != nil {
			return nil, err
		}

		objs = append(objs, &obj)
	}

	return objs, nil
}

func ModuleKubeConfigToDeckhouseConfig(kubeConfig *kcm.ModuleKubeConfig) (*d8config_v1.DeckhouseConfig, error) {
	return kubeConfigToDeckhouseConfig(
		kubeConfig.ModuleName,
		kubeConfig.ModuleConfigKey,
		kubeConfig.Values,
		kubeConfig.IsEnabled,
	)
}

func GlobalKubeConfigToDeckhouseConfig(kubeConfig *kcm.GlobalKubeConfig) (*d8config_v1.DeckhouseConfig, error) {
	return kubeConfigToDeckhouseConfig(
		"global",
		"global",
		kubeConfig.Values,
		nil,
	)
}

func kubeConfigToDeckhouseConfig(name string, valuesKey string, values utils.Values, isEnabled *bool) (*d8config_v1.DeckhouseConfig, error) {
	modConfig := DeckhouseConfigV0(name)

	// Fill the 'enabled' field.
	modConfig.Spec.Enabled = isEnabled

	untypedValues := values[valuesKey]
	incompatible := false
	switch v := untypedValues.(type) {
	case map[string]interface{}:
		modConfig.Spec.ConfigValues = v
	case nil:
		// Values are nil when only Enabled flag is present.
		// Handle this as empty 'configValues' field.
	case string:
		// Handle empty string as empty 'configValues' field.
		if v != "" {
			incompatible = true
		}
	case []interface{}:
		// Handle empty array as empty 'configValues' field.
		if len(v) != 0 {
			incompatible = true
		}
	default:
		incompatible = true
	}

	if incompatible {
		return nil, fmt.Errorf("configmap section '%s' is not an object, need map[string]interface{}, got %T:(%+v)", valuesKey, untypedValues, untypedValues)
	}
	return modConfig, nil
}

func DeckhouseConfigCR(name string) *d8config_v1.DeckhouseConfig {
	return &d8config_v1.DeckhouseConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeckhouseConfig",
			APIVersion: "deckhouse.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: d8config_v1.DeckhouseConfigSpec{
			ConfigVersion: "v1.0.0",
		},
	}
}

// DeckhouseConfigV0 returns DeckhouseConfig for values from cm/deckhouse.
func DeckhouseConfigV0(name string) *d8config_v1.DeckhouseConfig {
	return &d8config_v1.DeckhouseConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeckhouseConfig",
			APIVersion: "deckhouse.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: d8config_v1.DeckhouseConfigSpec{
			ConfigVersion: "v0.0.0",
		},
	}
}

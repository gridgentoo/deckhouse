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
	"fmt"
	"strconv"

	kcm "github.com/flant/addon-operator/pkg/kube_config_manager"
	"github.com/flant/addon-operator/pkg/utils"
	"sigs.k8s.io/yaml"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/modules"
	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
)

// RecreateDeckhouseConfigs returns DeckhouseConfig objects:
// - if new module appears in modules directory
// - if module section is present in generated ConfigMap but corresponding DeckhouseConfig is absent.
func RecreateDeckhouseConfigs(cmData map[string]string, allConfigs []*d8config_v1.DeckhouseConfig) ([]*d8config_v1.DeckhouseConfig, error) {
	kubeConfig, err := kcm.ParseConfigMapData(cmData)
	if err != nil {
		return nil, fmt.Errorf("parse cm/deckhouse data: %v", err)
	}

	knownModules := modules.Registry().GetModules()

	configNames := map[string]struct{}{}
	for _, cfg := range allConfigs {
		configNames[cfg.GetName()] = struct{}{}
	}

	newConfigs := make([]*d8config_v1.DeckhouseConfig, 0)

	// Try to re-create first.
	if kubeConfig.Global != nil {
		if _, has := configNames["global"]; !has {
			cfgObj, err := GlobalKubeConfigToDeckhouseConfig(kubeConfig.Global)
			if err != nil {
				return nil, err
			}
			newConfigs = append(newConfigs, cfgObj)
		}
	}
	for _, kubeCfg := range kubeConfig.Modules {
		if _, has := configNames[kubeCfg.ModuleName]; !has {
			cfgObj, err := ModuleKubeConfigToDeckhouseConfig(kubeCfg)
			if err != nil {
				return nil, err
			}
			newConfigs = append(newConfigs, cfgObj)
		}
	}

	// Create DeckhouseConfig objects for new modules.
	for _, mod := range knownModules {
		// Do not create empty DeckhouseConfig if ConfigMap has values for module.
		if _, has := kubeConfig.Modules[mod.Name]; has {
			continue
		}
		// Do not create empty DeckhouseConfig if DeckhouseConfig exists.
		if _, has := configNames[mod.Name]; has {
			continue
		}
		cfgObj := DeckhouseConfigCR(mod.Name)
		newConfigs = append(newConfigs, cfgObj)
	}

	// Run conversions for recreated and new modules.
	// No conversions for configs already in cluster.
	for _, cfg := range newConfigs {
		err := ValidateValues(cfg)
		if err != nil {
			return nil, err
		}
	}

	return newConfigs, nil
}

// SyncFromDeckhouseConfigs creates new Data for ConfigMap from DeckhouseConfig objects.
func SyncFromDeckhouseConfigs(allConfigs []*d8config_v1.DeckhouseConfig) (map[string]string, error) {
	data := make(map[string]string)

	for _, cfg := range allConfigs {
		name := cfg.GetName()

		valuesKey := utils.ModuleNameToValuesKey(name)
		if cfg.Spec.ConfigValues != nil {
			sectionBytes, err := yaml.Marshal(cfg.Spec.ConfigValues)
			if err != nil {
				return nil, err
			}
			data[valuesKey] = string(sectionBytes)
		}

		// Prevent creating 'globalEnabled' key.
		if name == "global" {
			continue
		}

		if cfg.Spec.Enabled != nil {
			enabledKey := valuesKey + "Enabled"
			data[enabledKey] = strconv.FormatBool(*cfg.Spec.Enabled)
		}
	}

	return data, nil
}

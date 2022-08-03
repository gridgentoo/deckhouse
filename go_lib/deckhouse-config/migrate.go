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

	kcm "github.com/flant/addon-operator/pkg/kube_config_manager"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	"github.com/deckhouse/deckhouse/go_lib/dependency/k8s"
	"github.com/deckhouse/deckhouse/go_lib/set"
)

// MigrateConfigMapToModuleConfigs return DeckhouseConfig object for global section
// and for each module specified in ConfigMap/deckhouse.
// It sets config version to v0.0.0 and config values to values from ConfigMap/deckhouse.
func MigrateConfigMapToModuleConfigs(kubeClient k8s.Client, possibleNames []string) ([]*d8config_v1.DeckhouseConfig, error) {
	// Split ConfigMap/deckhouse to ModuleSettings resources.
	deckhouseCm, err := GetDeckhouseConfigMap(kubeClient)
	if err != nil {
		return nil, fmt.Errorf("get cm/deckhouse: %v", err)
	}

	cfg, err := kcm.ParseConfigMapData(deckhouseCm.Data)
	if err != nil {
		return nil, fmt.Errorf("parse cm/deckhouse data: %v", err)
	}

	objs := make([]*d8config_v1.DeckhouseConfig, 0)

	// Convert global section to ModuleSettings.
	globalConfig, err := GlobalKubeConfigToDeckhouseConfig(cfg.Global)
	if err != nil {
		return nil, err
	}
	objs = append(objs, globalConfig)

	// Convert modules sections to ModuleSettings.
	// Note: possibleNames are kebab-cased, cfg.Module keys are also kebab-cased.
	possibleNamesSet := set.New(possibleNames...)
	for name, modCfg := range cfg.Modules {
		// Ignore migrating for unknown module names.
		if !possibleNamesSet.Has(name) {
			continue
		}
		modConfig, err := ModuleKubeConfigToDeckhouseConfig(modCfg)
		if err != nil {
			return nil, err
		}
		// If conversions are available for module, use them to convert to the latest version of values.
		if conversion.Registry().HasModule(modConfig.GetName()) {
			newVersion, newValues, err := conversion.ConvertToLatest(modConfig.GetName(), modConfig.Spec.ConfigVersion, modConfig.Spec.ConfigValues)
			if err != nil {
				return nil, err
			}
			modConfig.Spec.ConfigVersion = newVersion
			modConfig.Spec.ConfigValues = newValues
		}
		objs = append(objs, modConfig)
	}

	return objs, nil
}

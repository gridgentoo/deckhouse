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

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/modules"
	v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
)

func ValidateValues(cfg *v1.DeckhouseConfig) error {
	// Ignore conversion and validation for empty values if module is not enabled explicitly.
	if len(cfg.Spec.ConfigValues) == 0 && (cfg.Spec.Enabled == nil || !*cfg.Spec.Enabled) {
		return nil
	}

	origVersion := cfg.Spec.ConfigVersion

	// Run registered conversions.
	convertedMsg := ""
	if conversion.Registry().HasModule(cfg.GetName()) {
		newVersion, newValues, err := conversion.ConvertToLatest(cfg.GetName(), cfg.Spec.ConfigVersion, cfg.Spec.ConfigValues)
		if err != nil {
			return fmt.Errorf("convert %s config values from version %s to latest: %v", cfg.GetName(), cfg.Spec.ConfigVersion, err)
		}
		cfg.Spec.ConfigVersion = newVersion
		cfg.Spec.ConfigValues = newValues
		convertedMsg = fmt.Sprintf(" converted to %s", newVersion)
	}

	err := modules.Registry().ValidateConfigValues(cfg.GetName(), cfg.Spec.ConfigValues)
	if err != nil {
		return fmt.Errorf("%s config values of version %s%s are not valid: %v", cfg.GetName(), origVersion, convertedMsg, err)
	}

	return nil
}

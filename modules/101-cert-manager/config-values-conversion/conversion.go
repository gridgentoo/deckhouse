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

package config_values_conversion

import (
	"strings"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
)

const moduleName = "cert-manager"

var _ = conversion.RegisterFunc(moduleName, "v0.0.0", func(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error) {
	newVals, err := convertV0ToV1(configValues)
	return "v1.0.0", newVals, err
})

var _ = conversion.RegisterFunc(moduleName, "v1.0.0", func(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error) {
	newVals, err := convertV1ToV2(configValues)
	return "v2.0.0", newVals, err
})

var _ = conversion.RegisterFunc(moduleName, "v2.0.0", func(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error) {
	return "v2.0.0", configValues, nil
})

func convertV1ToV2(values map[string]interface{}) (map[string]interface{}, error) {
	return values, nil
}

// convertV0ToV1 migrate flat settings for cloudflare to an object.
func convertV0ToV1(values map[string]interface{}) (map[string]interface{}, error) {
	newValues := make(map[string]interface{})

	newCloudflare := make(map[string]interface{})
	for k, v := range values {
		if strings.HasPrefix(k, "cloudflare") {
			newKey := strings.TrimPrefix(k, "cloudflare")
			newCloudflare[newKey] = v
		} else {
			newValues[k] = v
		}
	}

	if len(newCloudflare) > 0 {
		newValues["cloudflare"] = newCloudflare
	}

	return newValues, nil
}

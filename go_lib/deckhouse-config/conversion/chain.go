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

package conversion

import "fmt"

func ConvertToLatest(moduleName string, fromVersion string, values map[string]interface{}) (string, map[string]interface{}, error) {
	maxTries := Registry().Count(moduleName)

	tries := 0
	version := fromVersion
	convValues := values
	for {
		conv := Registry().Get(moduleName, version)
		if conv == nil {
			return "", nil, fmt.Errorf("convert from %s: conversion chain interrupt: no conversion for %s", fromVersion, version)
		}
		newVer, newValues, err := conv.Convert(version, convValues)
		if err != nil {
			return "", nil, fmt.Errorf("convert from %s: conversion chain error for %s: %v", fromVersion, version, err)
		}

		version = newVer
		convValues = newValues

		// Conversion for the latest version returns same version.
		if version == newVer {
			return version, newValues, nil
		}

		// Prevent looped conversions.
		tries++
		if tries > maxTries {
			return "", nil, fmt.Errorf("convert from %s: conversion chain too long or looped", fromVersion)
		}
	}
}

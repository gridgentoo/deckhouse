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

type ConversionFunc func(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error)

type Conversion interface {
	Convert(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error)
}

type anonymousConversion struct {
	fn ConversionFunc
}

func (a *anonymousConversion) Convert(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error) {
	if a.fn != nil {
		return a.fn(configVersion, configValues)
	}
	return "", nil, nil
}

func NewAnonymousConversion(conversionFunc ConversionFunc) Conversion {
	return &anonymousConversion{fn: conversionFunc}
}

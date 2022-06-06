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

import "sync"

var instance *ConvRegistry
var once sync.Once

func Registry() *ConvRegistry {
	once.Do(func() {
		instance = new(ConvRegistry)
	})
	return instance
}

// Register adds Conversion implementation to Registry. Returns true to use with "var _ =".
func Register(moduleName string, srcVersion string, conversion Conversion) bool {
	Registry().Add(moduleName, srcVersion, conversion)
	return true
}

// RegisterFunc adds a function as a Conversion to Registry. Returns true to use with "var _ =".
func RegisterFunc(moduleName string, srcVersion string, conversionFunc ConversionFunc) bool {
	Registry().Add(moduleName, srcVersion, NewAnonymousConversion(conversionFunc))
	return true
}

type ConvRegistry struct {
	// module name -> version -> convertor
	conversions map[string]map[string]Conversion

	m sync.Mutex
}

func (r *ConvRegistry) Add(moduleName string, srcVersion string, conversion Conversion) {
	r.m.Lock()
	defer r.m.Unlock()

	if r.conversions == nil {
		r.conversions = make(map[string]map[string]Conversion)
	}
	if _, has := r.conversions[moduleName]; !has {
		r.conversions[moduleName] = make(map[string]Conversion)
	}

	r.conversions[moduleName][srcVersion] = conversion
}

func (r *ConvRegistry) Get(moduleName string, srcVersion string) Conversion {
	r.m.Lock()
	defer r.m.Unlock()
	if r.conversions == nil {
		return nil
	}
	if _, has := r.conversions[moduleName]; !has {
		return nil
	}
	return r.conversions[moduleName][srcVersion]
}

// Count returns a number of registered conversions for the module.
func (r *ConvRegistry) Count(moduleName string) int {
	r.m.Lock()
	defer r.m.Unlock()
	return len(r.conversions[moduleName])
}

func IsRegistered(moduleName string) bool {
	return Registry().Count(moduleName) > 0
}

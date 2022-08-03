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

import (
	"sync"

	"github.com/Masterminds/semver/v3"
)

/*
Conversion package is used to support older values layouts.

Module should define conversion functions and register them in conversion
Registry. Conversion webhook will use these functions to convert values in
DeckhouseConfig objects to latest version.
*/

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

	// module name -> latest version
	latestVersions map[string]*semver.Version

	m sync.RWMutex
}

func (r *ConvRegistry) Add(moduleName string, srcVersion string, conversion Conversion) {
	r.m.Lock()
	defer r.m.Unlock()

	srcSemver := semver.MustParse(srcVersion)

	if r.conversions == nil {
		r.conversions = make(map[string]map[string]Conversion)
	}
	if _, has := r.conversions[moduleName]; !has {
		r.conversions[moduleName] = make(map[string]Conversion)
	}

	r.conversions[moduleName][srcVersion] = conversion

	// Update latest version.
	if r.latestVersions == nil {
		r.latestVersions = make(map[string]*semver.Version)
	}

	storedVersion, has := r.latestVersions[moduleName]
	if !has || srcSemver.GreaterThan(storedVersion) {
		r.latestVersions[moduleName] = srcSemver
	}
}

func (r *ConvRegistry) Get(moduleName string, srcVersion string) Conversion {
	r.m.RLock()
	defer r.m.RUnlock()
	if r.conversions == nil {
		return nil
	}
	if _, has := r.conversions[moduleName]; !has {
		return nil
	}
	return r.conversions[moduleName][srcVersion]
}

func (r *ConvRegistry) LatestVersion(moduleName string) string {
	storedVersion := r.latestVersions[moduleName]
	if storedVersion != nil {
		return storedVersion.Original()
	}
	return ""
}

// Count returns a number of registered conversions for the module.
func (r *ConvRegistry) Count(moduleName string) int {
	r.m.RLock()
	defer r.m.RUnlock()
	return len(r.conversions[moduleName])
}

// HasModule returns whether module has registered conversions.
func (r *ConvRegistry) HasModule(moduleName string) bool {
	r.m.RLock()
	defer r.m.RUnlock()
	_, has := r.conversions[moduleName]
	return has
}

// HasVersion returns whether module has registered conversion for version.
func (r *ConvRegistry) HasVersion(moduleName string, version string) bool {
	r.m.RLock()
	defer r.m.RUnlock()
	_, has := r.conversions[moduleName]
	if has {
		_, has = r.conversions[moduleName][version]
		return has
	}
	return false
}

// VersionList returns all versions for the module.
func (r *ConvRegistry) VersionList(moduleName string) []string {
	r.m.RLock()
	defer r.m.RUnlock()
	versions := make([]string, 0)
	for ver := range r.conversions[moduleName] {
		versions = append(versions, ver)
	}
	return versions
}

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

package modules

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/flant/addon-operator/pkg/module_manager"
	"github.com/flant/addon-operator/pkg/values/validation"
)

const DefaultModulesDir = "/deckhouse/modules"
const DefaultGlobalHooksDir = "/deckhouse/global-hooks"

type ModulesRegistry interface {
	Init(globalHooksDir string, modulesDir string) error
	GetModules() []*module_manager.Module
	GetModule(cfgName string) *module_manager.Module
	HasModule(cfgName string) bool
	GetPossibleNames() []string
	ValidateConfigValues(cfgName string, cfgValues map[string]interface{}) error
}

type modulesRegistry struct {
	names           []string
	modules         []*module_manager.Module
	moduleMap       map[string]*module_manager.Module
	valuesValidator *validation.ValuesValidator
}

var instance *modulesRegistry
var once sync.Once

func Registry() ModulesRegistry {
	once.Do(func() {
		instance = &modulesRegistry{}
	})
	return instance
}

func NewModulesRegistry() ModulesRegistry {
	return &modulesRegistry{}
}

// Init searches for all modules in module directory and load
// OpenAPI schemas for global config values and module config values.
func (m *modulesRegistry) Init(globalHooksDir string, modulesDir string) (err error) {
	m.names = make([]string, 0)
	m.modules = make([]*module_manager.Module, 0)
	m.moduleMap = make(map[string]*module_manager.Module)
	m.valuesValidator = validation.NewValuesValidator()

	// Load OpenAPI schema for global config values.
	m.names = append(m.names, "global")
	openAPIPath := filepath.Join(globalHooksDir, "openapi")
	configBytes, _, err := module_manager.ReadOpenAPIFiles(openAPIPath)
	if err != nil {
		return fmt.Errorf("read global openAPI schemas: %v", err)
	}

	err = m.valuesValidator.SchemaStorage.AddGlobalValuesSchemas(
		configBytes,
		nil,
	)
	if err != nil {
		return fmt.Errorf("add global config values OpenAPI schema: %v", err)
	}

	m.modules, err = module_manager.SearchModules(modulesDir)
	if err != nil {
		return fmt.Errorf("load modules from %s: %v", modulesDir, err)
	}

	for _, module := range m.modules {
		m.moduleMap[module.Name] = module
		m.names = append(m.names, module.Name)
		// Load OpenAPI schema for module config values.
		openAPIPath := filepath.Join(module.Path, "openapi")
		configBytes, _, err := module_manager.ReadOpenAPIFiles(openAPIPath)
		if err != nil {
			return fmt.Errorf("module '%s' read OpenAPI schemas: %v", module.Name, err)
		}

		err = m.valuesValidator.SchemaStorage.AddModuleValuesSchemas(
			module.ValuesKey(),
			configBytes,
			nil,
		)
		if err != nil {
			return fmt.Errorf("add module '%s' config values OpenAPI schema: %v", module.Name, err)
		}
	}
	return nil
}

// GetModules returns modules in /deckhouse/modules directory
// with initialized openapi.
// Reduced module_manager.RegisterModules().
func (m *modulesRegistry) GetModules() []*module_manager.Module {
	return m.modules
}

// GetModule returns a module by name.
func (m *modulesRegistry) GetModule(cfgName string) *module_manager.Module {
	return m.moduleMap[cfgName]
}

// HasModule returns true if module is known.
func (m *modulesRegistry) HasModule(cfgName string) bool {
	_, has := m.moduleMap[cfgName]
	return has
}

// GetPossibleNames returns kebab-cased names for all modules
// plus "global".
func (m *modulesRegistry) GetPossibleNames() []string {
	return m.names
}

// ValidateConfigValues uses OpenAPI schema to validate input config values.
func (m *modulesRegistry) ValidateConfigValues(cfgName string, cfgValues map[string]interface{}) error {
	if cfgName == "global" {
		values := map[string]interface{}{
			"global": cfgValues,
		}

		return m.valuesValidator.ValidateGlobalConfigValues(values)
	}

	module := m.moduleMap[cfgName]
	valuesKey := module.ValuesKey()
	moduleValues := map[string]interface{}{
		valuesKey: cfgValues,
	}

	return m.valuesValidator.ValidateModuleConfigValues(valuesKey, moduleValues)
}

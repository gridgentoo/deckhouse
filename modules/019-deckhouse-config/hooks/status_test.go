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

package hooks

import (
	addon_operator "github.com/flant/addon-operator/pkg/addon-operator"
	"github.com/flant/addon-operator/pkg/module_manager"
	"github.com/flant/addon-operator/pkg/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var prometheusConfigYaml = `
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: prometheus
spec:
  configVersion: v1.0.0
  configValues:
    param1: val1
`

var enabledPrometheusConfigYaml = `
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: prometheus
spec:
  configVersion: v1.0.0
  configValues:
    param1: val1
  enabled: true
`

var _ = Describe("Modules :: deckhouse-config :: hooks :: update DeckhouseConfig status ::", func() {
	f := HookExecutionConfigInit(`{}`, `{}`)

	f.RegisterCRD("deckhouse.io", "v1", "DeckhouseConfig", false)

	Context("Known module enabled", func() {
		BeforeEach(func() {
			f.ValuesSetFromYaml(PossibleNamesPath, []byte(`["prometheus", "cert-manager"]`))
			f.KubeStateSet(prometheusConfigYaml)

			dependency.TestDC.AddonOperator = MockAddonOperator(map[string]*module_manager.Module{
				"prometheus": newMockedModule("prometheus", nil, nil),
			}, map[string]struct{}{
				"prometheus": {},
			})

			f.BindingContexts.Set(f.GenerateScheduleContext("*/15 * * * * *"))
			f.RunHook()
		})

		It("Should be enabled in status", func() {
			Expect(f).To(ExecuteSuccessfully())

			promCfg := f.KubernetesGlobalResource("DeckhouseConfig", "prometheus")
			Expect(promCfg.Field("status.status").String()).To(Equal("Enabled"), "should update status")
		})
	})

	Context("Enabled by bundle, disabled by enabled script", func() {
		BeforeEach(func() {
			f.ValuesSetFromYaml(PossibleNamesPath, []byte(`["prometheus", "cert-manager"]`))
			f.KubeStateSet(enabledPrometheusConfigYaml)

			dependency.TestDC.AddonOperator = MockAddonOperator(map[string]*module_manager.Module{
				"prometheus": newMockedModule("prometheus", pointer.Bool(true), nil),
			}, map[string]struct{}{})

			f.BindingContexts.Set(f.GenerateScheduleContext("*/15 * * * * *"))
			f.RunHook()
		})

		It("Should not be enabled in status", func() {
			Expect(f).To(ExecuteSuccessfully())

			promCfg := f.KubernetesGlobalResource("DeckhouseConfig", "prometheus")
			Expect(promCfg.Field("status.status").String()).To(ContainSubstring("Disabled by script"), "should update status")
		})
	})
})

// MockAddonOperator returns AddonOperator instance suitable for testing
// status.go hook.
func MockAddonOperator(modules map[string]*module_manager.Module, enabledModules map[string]struct{}) *addon_operator.AddonOperator {
	return &addon_operator.AddonOperator{
		ModuleManager: &ModuleManagerMock{
			modules:        modules,
			enabledModules: enabledModules,
		},
	}
}

type ModuleManagerMock struct {
	module_manager.ModuleManager
	modules        map[string]*module_manager.Module
	enabledModules map[string]struct{}
}

func (m *ModuleManagerMock) IsModuleEnabled(name string) bool {
	_, has := m.enabledModules[name]
	return has
}

func (m *ModuleManagerMock) GetModule(name string) *module_manager.Module {
	mod, has := m.modules[name]
	if has {
		return mod
	}
	return nil
}

func newMockedModule(name string, commonEnabled *bool, staticEnabled *bool) *module_manager.Module {
	return &module_manager.Module{
		Name: name,
		CommonStaticConfig: &utils.ModuleConfig{
			IsEnabled: commonEnabled,
		},
		StaticConfig: &utils.ModuleConfig{
			IsEnabled: staticEnabled,
		},
	}
}

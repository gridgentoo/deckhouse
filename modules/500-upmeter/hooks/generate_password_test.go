/*
Copyright 2021 Flant JSC

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
	"encoding/base64"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

type appSettings struct {
	password                   string
	passwordB64                string
	nsManifest                 string
	secretManifest             string
	externalAuthValuesPath     string
	passwordValuesPath         string
	passwordInternalValuesPath string
}

var _ = Describe("Modules :: upmeter :: hooks :: generate_password", func() {
	testSettings := make(map[string]*appSettings)
	for secretName, appName := range upmeterApps {
		settings := new(appSettings)

		settings.password = fmt.Sprintf("t3stPassw0rd-%s", appName)
		settings.passwordB64 = base64.StdEncoding.EncodeToString([]byte(settings.password))
		settings.nsManifest = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: ` + upmeterNS + "\n"

		// Secret with password.
		settings.secretManifest = `
---
apiVersion: v1
kind: Secret
metadata:
  name: ` + secretName + `
  namespace: ` + upmeterNS + `
data:
  auth: ` + settings.passwordB64 + "\n"

		settings.externalAuthValuesPath = fmt.Sprintf(externalAuthValuesTmpl, appName)
		settings.passwordValuesPath = fmt.Sprintf(passwordValuesTmpl, appName)
		settings.passwordInternalValuesPath = fmt.Sprintf(passwordInternalValuesTmpl, appName)

		testSettings[appName] = settings
	}

	for appName, settings := range testSettings {
		Context(appName, func() {

			// Initialize internal.auth object for values patch to work.
			f := HookExecutionConfigInit(
				`{"upmeter": {"internal": {"auth": {"status": {}, "webui": {}}}} }`,
				`{"upmeter":{}}`,
			)

			Context("with no auth settings", func() {
				BeforeEach(func() {
					f.KubeStateSet("")
					f.BindingContexts.Set(f.GenerateBeforeHelmContext())
					f.RunHook()
				})

				It("should generate new password", func() {
					Expect(f).To(ExecuteSuccessfully())
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).ShouldNot(BeEmpty())
				})
			})

			Context("with password in configuration", func() {
				BeforeEach(func() {
					f.KubeStateSet("")
					f.BindingContexts.Set(f.GenerateBeforeHelmContext())
					f.ConfigValuesSet(settings.passwordValuesPath, settings.password)
					f.RunHook()
				})

				It("should set password from configuration", func() {
					Expect(f).To(ExecuteSuccessfully())
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).Should(BeEquivalentTo(settings.password))
				})
			})

			Context("with external auth", func() {
				BeforeEach(func() {
					f.KubeStateSet("")
					f.BindingContexts.Set(f.GenerateBeforeHelmContext())
					f.ValuesSetFromYaml(settings.externalAuthValuesPath, []byte(`{"authURL": "test"}`))
					f.RunHook()
				})

				It("should clean password from values", func() {
					Expect(f).To(ExecuteSuccessfully())
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).Should(BeEmpty())
				})
			})

			Context("with password in Secret", func() {
				BeforeEach(func() {
					f.KubeStateSet(settings.nsManifest + settings.secretManifest)
					f.BindingContexts.Set(f.GenerateBeforeHelmContext())
					f.ValuesSet(settings.passwordValuesPath, "not-a-test-password")
					f.RunHook()
				})
				It("should set password from Secret", func() {
					Expect(f).To(ExecuteSuccessfully())
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).Should(BeEquivalentTo(settings.password))
				})
			})

		})

	}

	Context("all apps", func() {
		f := HookExecutionConfigInit(
			`{"upmeter": {"internal": {"auth": {"status": {}, "webui": {}}}} }`,
			`{"upmeter":{}}`,
		)

		Context("with no auth settings", func() {
			BeforeEach(func() {
				f.KubeStateSet("")
				f.BindingContexts.Set(f.GenerateBeforeHelmContext())
				f.RunHook()
			})

			It("should generate new password", func() {
				Expect(f).To(ExecuteSuccessfully())

				for appName, settings := range testSettings {
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).ShouldNot(BeEmpty(), "Should generate password for '%s'", appName)
				}
			})
		})

		Context("with passwords in configuration", func() {

			BeforeEach(func() {
				f.KubeStateSet("")
				f.BindingContexts.Set(f.GenerateBeforeHelmContext())

				for _, settings := range testSettings {
					f.ConfigValuesSet(settings.passwordValuesPath, settings.password)
				}

				f.RunHook()
			})

			It("should set password from configuration", func() {
				Expect(f).To(ExecuteSuccessfully())
				for appName, settings := range testSettings {
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).Should(Equal(settings.password), "Should set password from configuration for '%s'", appName)
				}
			})
		})

		Context("with external auth", func() {
			BeforeEach(func() {
				f.KubeStateSet("")
				f.BindingContexts.Set(f.GenerateBeforeHelmContext())

				extAuth := []byte(`{"authURL": "test"}`)

				for _, settings := range testSettings {
					f.ValuesSetFromYaml(settings.externalAuthValuesPath, extAuth)
				}
				f.RunHook()
			})

			It("should clean password from values", func() {
				Expect(f).To(ExecuteSuccessfully())

				for _, settings := range testSettings {
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).Should(BeEmpty())
				}
			})
		})

		Context("with passwords in Secrets", func() {
			BeforeEach(func() {
				kubeState := ""
				for _, settings := range testSettings {
					if kubeState == "" {
						kubeState = settings.nsManifest
					}
					kubeState += settings.secretManifest
					f.ValuesSet(settings.passwordValuesPath, "not-a-test-password")
				}
				f.KubeStateSet(kubeState)
				f.BindingContexts.Set(f.GenerateBeforeHelmContext())
				f.RunHook()
			})
			It("should set password from Secret", func() {
				Expect(f).To(ExecuteSuccessfully())
				for _, settings := range testSettings {
					Expect(f.ValuesGet(settings.passwordInternalValuesPath).String()).Should(BeEquivalentTo(settings.password))
				}
			})
		})

	})
})

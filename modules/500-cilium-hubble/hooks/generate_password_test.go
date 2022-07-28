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

	"github.com/deckhouse/deckhouse/go_lib/hooks/generate_password"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Modules :: cilium-hubble :: hooks :: generate_password", func() {
	var (
		hook = generate_password.NewBasicAuthPlainHook(moduleValuesKey, authSecretNS, authSecretName)

		testPassword    = "t3stPassw0rd"
		testPasswordB64 = base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("admin:{PLAIN}%s", testPassword),
		))

		// Namespace should be created before creating the Secret.
		nsManifest = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: ` + authSecretNS + "\n"

		// Secret with password.
		authSecretManifest = `
---
apiVersion: v1
kind: Secret
metadata:
  name: ` + authSecretName + `
  namespace: ` + authSecretNS + `
data:
  auth: ` + testPasswordB64 + "\n"
	)

	// Initialize internal.auth object for values patch to work.
	f := HookExecutionConfigInit(
		`{"ciliumHubble": {"internal": {"auth":{}}} }`,
		`{"ciliumHubble":{}}`,
	)

	Context("with no auth settings", func() {
		BeforeEach(func() {
			f.KubeStateSet("")
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})
		It("should generate new password", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet(hook.PasswordInternalKey()).String()).ShouldNot(BeEmpty())
		})
	})

	Context("with password in configuration", func() {
		BeforeEach(func() {
			f.KubeStateSet("")
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.ConfigValuesSet(hook.PasswordKey(), testPassword)
			f.RunHook()
		})
		It("should set password from configuration", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet(hook.PasswordInternalKey()).String()).Should(BeEquivalentTo(testPassword))
		})
	})

	Context("with external auth", func() {
		BeforeEach(func() {
			f.KubeStateSet("")
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.ValuesSetFromYaml(hook.ExternalAuthKey(), []byte(`{"authURL": "test"}`))
			f.RunHook()
		})
		It("should clean password from values", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet(hook.PasswordKey()).String()).Should(BeEmpty())
			Expect(f.ValuesGet(hook.PasswordInternalKey()).Exists()).Should(BeFalse())
		})
	})

	Context("with password in Secret", func() {
		BeforeEach(func() {
			f.KubeStateSet(nsManifest + authSecretManifest)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.ValuesSet(hook.PasswordKey(), "not-a-test-password")
			f.RunHook()
		})
		It("should set password from Secret", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet(hook.PasswordInternalKey()).String()).Should(BeEquivalentTo(testPassword))
		})
	})

})

/*
Copyright 2021 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
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

var _ = Describe("Modules :: istio :: hooks :: generate_password ", func() {
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

	f := HookExecutionConfigInit(
		`{"istio":{"internal":{"auth": {}}}}`,
		`{"istio":{}}`,
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

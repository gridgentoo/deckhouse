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

package main

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func runValidate(t *testing.T, rootDir string, manifest string) (*kwhvalidating.ValidatorResult, error) {
	g := NewWithT(t)

	cfgValidator := NewDeckhouseConfigValidator()
	cfgValidator.ModulesDir = filepath.Join(rootDir, "modules")
	cfgValidator.GlobalHooksDir = filepath.Join(rootDir, "global-hooks")

	err := cfgValidator.Init()
	g.Expect(err).ShouldNot(HaveOccurred())

	obj := modCfgFromManifest(manifest)

	return cfgValidator.Validate(context.TODO(), nil, obj)
}

const validModuleCfg = `
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: module-one
spec:
  configVersion: v1.0.0
  configValues:
    param1: someText
`

func Test_moduleConfigValidate_valid_object(t *testing.T) {
	g := NewWithT(t)

	res, err := runValidate(t, "testdata", validModuleCfg)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(res.Valid).Should(BeTrue())
	g.Expect(res.Warnings).Should(HaveLen(0))
}

const unknownModuleCfg = `
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: module-three
spec:
  configVersion: v1.0.0
  configValues:
    param1: someText
`

func Test_moduleConfigValidate_unknown_module(t *testing.T) {
	g := NewWithT(t)

	res, err := runValidate(t, "testdata", unknownModuleCfg)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(res.Valid).Should(BeTrue())
	g.Expect(res.Warnings).Should(HaveLen(1))
	g.Expect(res.Warnings[0]).Should(ContainSubstring("unknown"))
}

const invalidModuleCfg = `
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: module-one
spec:
  configVersion: v1.0.0
  configValues:
    param-forbidden: someText
`

func Test_moduleConfigValidate_invalid_object(t *testing.T) {
	g := NewWithT(t)

	res, err := runValidate(t, "testdata", invalidModuleCfg)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(res.Valid).Should(BeFalse())
	g.Expect(res.Message).Should(ContainSubstring("not valid"))
}

const validGlobalCfg = `
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: global
spec:
  configVersion: v1.0.0
  configValues:
    globalParam: someText
`

func Test_moduleConfigValidate_valid_global_object(t *testing.T) {
	g := NewWithT(t)

	res, err := runValidate(t, "testdata", validGlobalCfg)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(res.Valid).Should(BeTrue())
	g.Expect(res.Warnings).Should(HaveLen(0))
}

const invalidGlobalCfg = `
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: global
spec:
  configVersion: v1.0.0
  configValues:
    globalBadParam: someText
`

func Test_moduleConfigValidate_invalid_global_object(t *testing.T) {
	g := NewWithT(t)

	res, err := runValidate(t, "testdata", invalidGlobalCfg)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(res.Valid).Should(BeFalse())
	g.Expect(res.Message).Should(ContainSubstring("not valid"))
}

const noVerGlobalCfg = `
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: global
spec:
  configValues:
    globalParam: someText
`

func Test_moduleConfigValidate_no_version_global_object(t *testing.T) {
	g := NewWithT(t)

	res, err := runValidate(t, "testdata", noVerGlobalCfg)
	g.Expect(err).ShouldNot(HaveOccurred())

	g.Expect(res.Valid).Should(BeFalse())
	g.Expect(res.Message).Should(ContainSubstring("required"))
}

func modCfgFromManifest(manifest string) *unstructured.Unstructured {
	var m map[string]interface{}
	_ = yaml.Unmarshal([]byte(manifest), &m)
	return &unstructured.Unstructured{
		Object: m,
	}
}

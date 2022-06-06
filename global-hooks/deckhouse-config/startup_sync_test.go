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
	"context"

	"github.com/flant/addon-operator/sdk"
	"github.com/flant/kube-client/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	d8config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/modules"
	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Global hooks :: deckhouse-config :: migrate", func() {
	f := HookExecutionConfigInit(`{"global": {"discovery": {}}}`, `{}`)
	// Emulate ensure_crd hook.
	f.RegisterCRD("deckhouse.io", "v1", "DeckhouseConfig", false)

	registryErr := modules.Registry().Init(modules.DefaultGlobalHooksDir, modules.DefaultModulesDir)

	Context("Migrate to DeckhouseConfig objects from global section only", func() {
		BeforeEach(func() {
			f.KubeStateSet(``)

			createDeckhouseDeployment(f.KubeClient(), nil)
			createDeckhouseConfigMap(f.KubeClient(), map[string]string{
				"global": `
param1: val1
param2: val2
`,
			})

			f.BindingContexts.Set(f.GenerateOnStartupContext())
			f.RunHook()
		})

		It("Should have global section", func() {
			Expect(registryErr).ShouldNot(HaveOccurred(), "modules registry should be inited")
			Expect(f).To(ExecuteSuccessfully())
			cfgObjs, err := f.KubeClient().Dynamic().Resource(d8config_v1.GroupVersionResource()).List(context.TODO(), metav1.ListOptions{})
			cfgList := expectConfigItems(cfgObjs, modules.Registry().GetPossibleNames(), err)

			// Test global section.
			globalCfg := cfgList["global"]
			Expect(globalCfg).ShouldNot(BeNil())
			Expect(globalCfg.Spec.ConfigValues).ShouldNot(BeEmpty())
			Expect(globalCfg.Spec.ConfigValues).Should(HaveKey("param1"))
			Expect(globalCfg.Spec.ConfigValues).Should(HaveKey("param2"))

			// Test deployment
			deckhouseDeploy := f.KubernetesResource("Deployment", "d8-system", "deckhouse")
			Expect(deckhouseDeploy.Exists()).Should(BeTrue())
			Expect(deckhouseDeploy.Field("spec.template.spec.containers.0.env.1.value").String()).Should(Equal(d8config.GeneratedConfigMapName), "should update deploy/deckhouse to use generated ConfigMap")
		})
	})

	// NOTE: This test uses existing modules names and non-valid parameters.
	// TODO Create some dumb modules in testdata.
	Context("Migrate to DeckhouseConfig objects from module sections", func() {
		BeforeEach(func() {
			f.KubeStateSet(``)

			createDeckhouseDeployment(f.KubeClient(), nil)
			createDeckhouseConfigMap(f.KubeClient(), map[string]string{
				"global": `
param1: val1
param2: val2
`,
				"deckhouse": `
p1: val1
param1: val1
p2: val2
`,
				"certManager": `
p1: val1
p2: val2
param1: val1
`,
			})

			f.BindingContexts.Set(f.GenerateOnStartupContext())
			f.RunHook()
		})

		It("Should have param1 in global, deckhouse and cert-manager", func() {
			Expect(registryErr).ShouldNot(HaveOccurred(), "modules registry should be inited")
			Expect(f).To(ExecuteSuccessfully())
			cfgObjs, err := f.KubeClient().Dynamic().Resource(d8config_v1.GroupVersionResource()).List(context.TODO(), metav1.ListOptions{})

			cfgList := expectConfigItems(cfgObjs, modules.Registry().GetPossibleNames(), err)

			for _, name := range []string{"global", "deckhouse", "cert-manager"} {
				cfg := cfgList[name]
				Expect(cfg).ShouldNot(BeNil(), "DeckhouseConfig/%s should not be nil", name)
				Expect(cfg.Spec.ConfigValues).ShouldNot(BeEmpty(), "DeckhouseConfig/%s should have configValues", name)
				Expect(cfg.Spec.ConfigValues["param1"]).Should(Equal("val1"), "DeckhouseConfig/%s should have param1=val1", name)
			}
		})
	})
})

// TODO add more tests with modules in testdata.
var _ = Describe("Global hooks :: deckhouse-config :: sync", func() {
	f := HookExecutionConfigInit(`{"global": {"discovery": {}}}`, `{}`)
	// Emulate ensure_crd hook.
	f.RegisterCRD("deckhouse.io", "v1", "DeckhouseConfig", false)

	registryErr := modules.Registry().Init(modules.DefaultGlobalHooksDir, modules.DefaultModulesDir)

	Context("Sync new modules, create absent ConfigMap", func() {
		BeforeEach(func() {
			existingModuleConfigs := `
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: global
spec:
  configVersion: v1.0.0
  configValues:
    param1: val1
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: deckhouse
spec:
  configVersion: v1.0.0
  configValues:
    logLevel: Debug
`

			f.KubeStateSet(existingModuleConfigs)

			f.BindingContexts.Set(f.GenerateOnStartupContext())
			f.RunHook()
		})

		It("Should create DeckhouseConfig for all modules", func() {
			Expect(registryErr).ShouldNot(HaveOccurred(), "modules registry should be inited")
			Expect(f).To(ExecuteSuccessfully())

			cfgList, err := f.KubeClient().Dynamic().Resource(d8config_v1.GroupVersionResource()).List(context.TODO(), metav1.ListOptions{})
			_ = expectConfigItems(cfgList, modules.Registry().GetPossibleNames(), err)

			gcm := f.KubernetesResource("ConfigMap", "d8-system", d8config.GeneratedConfigMapName)
			Expect(gcm.Exists()).Should(BeTrue(), "should create ConfigMap from DeckhouseConfig")
			Expect(gcm.Field("data").Map()).Should(HaveLen(2), "generated ConfigMap should have sections only for existing DeckhouseConfig objects")
			Expect(gcm.Field("data.global").String()).Should(Equal("param1: val1\n"))
			Expect(gcm.Field("data.deckhouse").String()).Should(Equal("logLevel: Debug\n"))
		})
	})

	Context("Sync new modules, update existing ConfigMap", func() {
		BeforeEach(func() {
			existingConfigs := `
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: global
spec:
  configVersion: v1.0.0
  configValues:
    param1: val1
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: deckhouse
spec:
  configVersion: v1.0.0
  configValues:
    logLevel: Debug
`

			f.KubeStateSet(existingConfigs)

			// Use typed clientset for ConfigMap as the hook does.
			cm := d8config.GeneratedConfigMap(map[string]string{
				"global": "param2: val4",
			})
			_, _ = f.KubeClient().CoreV1().ConfigMaps("d8-system").Create(context.TODO(), cm, metav1.CreateOptions{})

			f.BindingContexts.Set(f.GenerateOnStartupContext())
			f.RunHook()
		})

		It("Should create DeckhouseConfig for all modules", func() {
			Expect(registryErr).ShouldNot(HaveOccurred(), "modules registry should be inited")
			Expect(f).To(ExecuteSuccessfully())

			cfgList, err := f.KubeClient().Dynamic().Resource(d8config_v1.GroupVersionResource()).List(context.TODO(), metav1.ListOptions{})
			_ = expectConfigItems(cfgList, modules.Registry().GetPossibleNames(), err)

			gcm := f.KubernetesResource("ConfigMap", "d8-system", d8config.GeneratedConfigMapName)
			Expect(gcm.Exists()).Should(BeTrue(), "should not delete generated ConfigMap")
			Expect(gcm.Field("data").Map()).Should(HaveLen(2), "generated ConfigMap should have sections only for existing DeckhouseConfig objects")
			Expect(gcm.Field("data.global").String()).Should(Equal("param1: val1\n"))
			Expect(gcm.Field("data.deckhouse").String()).Should(Equal("logLevel: Debug\n"))
		})
	})

	Context("Sync new modules: update existing ConfigMap", func() {
		BeforeEach(func() {
			existingModuleConfigs := `
---
apiVersion: deckhouse.io/v1
kind: DeckhouseConfig
metadata:
  name: global
spec:
  configVersion: v1.0.0
  configValues:
    param1: val1
`
			f.KubeStateSet(existingModuleConfigs)

			// Use typed clientset for ConfigMap as the hook does.
			cm := d8config.GeneratedConfigMap(map[string]string{
				"global":    "param2: val4",
				"deckhouse": "logLevel: Debug",
			})
			_, _ = f.KubeClient().CoreV1().ConfigMaps("d8-system").Create(context.TODO(), cm, metav1.CreateOptions{})

			f.BindingContexts.Set(f.GenerateOnStartupContext())
			f.RunHook()
		})

		It("Should create DeckhouseConfig for all modules", func() {
			Expect(registryErr).ShouldNot(HaveOccurred(), "modules registry should be inited")
			Expect(f).To(ExecuteSuccessfully())

			cfgList, err := f.KubeClient().Dynamic().Resource(d8config_v1.GroupVersionResource()).List(context.TODO(), metav1.ListOptions{})
			_ = expectConfigItems(cfgList, modules.Registry().GetPossibleNames(), err)

			gcm := f.KubernetesResource("ConfigMap", "d8-system", d8config.GeneratedConfigMapName)
			Expect(gcm.Exists()).Should(BeTrue(), "should not delete generated ConfigMap")
			Expect(gcm.Field("data").Map()).Should(HaveLen(1), "generated ConfigMap should have sections only for existing DeckhouseConfig objects")
			Expect(gcm.Field("data.global").String()).Should(Equal("param1: val1\n"))
			//Expect(gcm.Field("data.deckhouse").String()).Should(Equal("logLevel: Debug\n"))

			//deckhouseCfg := f.KubernetesGlobalResource("DeckhouseConfig", "deckhouse")
			//Expect(deckhouseCfg.Exists()).Should(BeTrue(), "should have DeckhouseConfig/deckhouse")
			//Expect(deckhouseCfg.Field("spec.configValues.logLevel").String()).Should(Equal("Debug"), "DeckhouseConfig/deckhouse should have logLevel=Debug from ConfigMap")
		})
	})
})

func expectConfigItems(list *unstructured.UnstructuredList, possibleNames []string, err error) map[string]*d8config_v1.DeckhouseConfig {
	Expect(err).ShouldNot(HaveOccurred(), "should get list of existing DeckhouseConfig objects")
	//Expect(list.Items).Should(HaveLen(len(possibleNames)), "should create DeckhouseConfig for each possible name")

	cfgs := make(map[string]*d8config_v1.DeckhouseConfig)
	for _, cfg := range list.Items {
		cfgObj := new(d8config_v1.DeckhouseConfig)
		err := sdk.FromUnstructured(&cfg, cfgObj)
		Expect(err).ShouldNot(HaveOccurred(), "should create valid DeckhouseConfig/%s object", cfg.GetName())
		cfgs[cfg.GetName()] = cfgObj
	}

	return cfgs
}

func deckhouseConfigMap(data map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deckhouse",
			Namespace: "d8-system",
		},
		Data: data,
	}
}

func createDeckhouseConfigMap(kubeClient client.Client, data map[string]string) {
	cm := deckhouseConfigMap(data)
	_, _ = kubeClient.CoreV1().ConfigMaps("d8-system").Create(context.TODO(), cm, metav1.CreateOptions{})
}

func deckhouseDeployment(envs map[string]string) *appsv1.Deployment {
	defaultEnvs := []string{
		"LOG_LEVEL", "Info",
		"ADDON_OPERATOR_CONFIG_MAP", "deckhouse",
		"HELM_HISTORY_MAX", "3",
	}
	deployEnvs := make([]v1.EnvVar, 0)
	for i := 0; i < len(defaultEnvs); i += 2 {
		name := defaultEnvs[i]
		val := defaultEnvs[i+1]
		if override, has := envs[name]; has {
			val = override
		}
		deployEnvs = append(deployEnvs, v1.EnvVar{
			Name:  name,
			Value: val,
		})
	}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deckhouse",
			Namespace: "d8-system",
		},
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "deckhouse",
							Env:  deployEnvs,
						},
					},
				},
			},
		},
	}
}

func createDeckhouseDeployment(kubeClient client.Client, envs map[string]string) {
	depl := deckhouseDeployment(envs)

	// Use Dynamic as hook use Filter to update an object.
	u, _ := sdk.ToUnstructured(depl)
	_, _ = kubeClient.Dynamic().Resource(schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}).Namespace("d8-system").Create(context.TODO(), u, metav1.CreateOptions{})
}

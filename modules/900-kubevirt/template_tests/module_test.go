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

package template_tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/helm"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

const (
	globalValues = `
  enabledModules: ["vertical-pod-autoscaler-crd"]
  highAvailability: true
  modules:
    placement: {}
  discovery:
    kubernetesVersion: 1.21.9
    d8SpecificNodeCountByRole:
      worker: 3
      master: 3
  modulesImages:
    registry: registry.deckhouse.io
    registryDockercfg: Y2ZnCg==
    tags:
      common:
        alpine: be0a73850874a7000b223b461e37c32263e95574d379d4ea3305006e-1624978147531
        csiExternalAttacher119: 9b875cd5df6c7c8a27d5ae6f6fa4b90e05797857c2a2d434d90fa384-1624991865774
        csiExternalAttacher120: 53055872b801779298431ef159edc41c3270059750e41ad7ce547815-1624991942422
        csiExternalAttacher121: f553f9fe2efa329f1df967dc1ad2433db463af49675cbabeddc2eb54-1644311199940
        csiExternalAttacher122: f553f9fe2efa329f1df967dc1ad2433db463af49675cbabeddc2eb54-1644311199940
        csiExternalProvisioner119: a71ffdcb6ea573ecc44d7cbe22ba713331e92e08c8768023ecd62be0-1624991501421
        csiExternalProvisioner120: 7c94147745b7f3b3bc85cf0466e4a9d47027064cb95e0a776638d88c-1624991909639
        csiExternalProvisioner121: 2a208c02726d0fe0dad5298b2185661280664f4180968293e897cebd-1644311194495
        csiExternalProvisioner122: 2a208c02726d0fe0dad5298b2185661280664f4180968293e897cebd-1644311194495
        csiExternalResizer119: 2267cc7565514b5968ec95c99c3aca19d1930d40e2e34845f7723774-1624992194047
        csiExternalResizer120: 217459a0507b191a1744974c564df5e7de5046e5711bbafcf341a978-1624992095833
        csiExternalResizer121: e21da1d5629d422962e8dc7fce05c2a0777bc53fb7088bb71bff91b8-1644311196688
        csiExternalResizer122: e21da1d5629d422962e8dc7fce05c2a0777bc53fb7088bb71bff91b8-1644311196688
        csiExternalSnapshotter121: d554e018513a74708418c139ba57aeb75f782209d86980a729a4d97a-1644311215309
        csiExternalSnapshotter122: d554e018513a74708418c139ba57aeb75f782209d86980a729a4d97a-1644311215309
        csiLivenessprobe121: cb71a0e242c71416ee53e4e02fe0785d47c4147ee4df8bdcf0d2dfcc-1644311213440
        csiLivenessprobe122: cb71a0e242c71416ee53e4e02fe0785d47c4147ee4df8bdcf0d2dfcc-1644311213440
        csiNodeDriverRegistrar119: ecd5587d4fa58d22f91609028527b03da2f7ab9ed3c190ef90179c4b-1624991634159
        csiNodeDriverRegistrar120: 29bdbc548dd7bcab6d9cbc6156afd70870fe6818e9ea48d3a9a7eccd-1624991527573
        csiNodeDriverRegistrar121: 121aea9e665ab86dc8a83e3198acbf6a76c894c1612003c808405fb6-1638785322286
        csiNodeDriverRegistrar122: 121aea9e665ab86dc8a83e3198acbf6a76c894c1612003c808405fb6-1638785322286
        kubeRbacProxy: a4506c2aa962611cf1858c774129d2a4f233502ecc376929aa97b9f5-1639403210069
        pause: 47e1a07baefaaa885a306bb24546c929b910fe9cffffd07218b66c0a-1624979682719
      kubevirt:
        virtHandler: e6499c8480a4a68e185da2ea03411743d607b1b4be9ad3846ef5394c-1660062436360
        virtOperator: 2de8b94f9f2c14331396df6f376c79b25d331af7a1fb3b406ff145f9-1660062522285
        virtApi: 7401840b4b0fe9506cb71b3603a97b2f6131d0d361829fb37901c592-1660062424656
        virtLauncher: f44a82e94956fa0c0600e258f9948cf6a4e39be5d309bd63b8864d7a-1660062320015
        virtController: 976bdbcb0648c544ca4653ddfdda071258780a6f64c5d33e7917c3ba-1660062417922
        vmiRouter: b8b7cc1b4610b59c498f517b7a03ca27455811b0efe8815bcefb624c-1660049501622
`
	moduleValues = ``
)

var _ = Describe("Module :: kubevirt :: helm template ::", func() {
	f := SetupHelmConfig(``)

	Context("Standard setup with SSL", func() {
		BeforeEach(func() {
			f.ValuesSetFromYaml("global", globalValues)
			f.ValuesSetFromYaml("kubevirt", moduleValues)
			f.HelmRender()
		})

		It("Everything must render properly", func() {
			Expect(f.RenderError).ShouldNot(HaveOccurred())
		})

	})
})

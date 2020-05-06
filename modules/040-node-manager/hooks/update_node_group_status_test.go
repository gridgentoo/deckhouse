package hooks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Modules :: node-manager :: hooks :: update_node_group_status ::", func() {
	const (
		stateCloudNG1 = `
---
apiVersion: deckhouse.io/v1alpha1
kind: NodeGroups
metadata:
  name: ng1
spec:
  nodeType: Cloud
  cloudInstances:
    maxPerZone: 5
    minPerZone: 1
status:
  extra: thing
`
		stateCloudNG2 = `
---
apiVersion: deckhouse.io/v1alpha1
kind: NodeGroups
metadata:
  name: ng2
spec:
  nodeType: Cloud
  cloudInstances:
    maxPerZone: 3
    minPerZone: 2
    zones: [a, b, c]
status: {}
`
		stateNG1 = `
---
apiVersion: deckhouse.io/v1alpha1
kind: NodeGroups
metadata:
  name: ng1
spec:
  nodeType: Static
status:
  extra: thing
`
		stateNG2 = `
---
apiVersion: deckhouse.io/v1alpha1
kind: NodeGroups
metadata:
  name: ng2
spec:
  nodeType: Static
status: {}
`
		stateMDs = `
---
apiVersion: machine.sapcloud.io/v1alpha1
kind: MachineDeployment
metadata:
  name: md-ng1
  namespace: d8-cloud-instance-manager
  labels:
    node-group: ng1
spec:
  replicas: 2
`
		stateMachines = `
---
apiVersion: machine.sapcloud.io/v1alpha1
kind: Machine
metadata:
  name: machine-ng1-aaa
  namespace: d8-cloud-instance-manager
  labels:
    instance-group: ng1-nova
---
apiVersion: machine.sapcloud.io/v1alpha1
kind: Machine
metadata:
  name: machine-ng1-bbb
  namespace: d8-cloud-instance-manager
  labels:
    instance-group: ng1-nova
`

		stateNodes = `
---
apiVersion: v1
kind: Node
metadata:
  name: node-ng1-aaa
  labels:
    node.deckhouse.io/group: ng1
  annotations:
    node.deckhouse.io/configuration-checksum: a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3
status:
  conditions:
  - some: thing
  - status: "False"
    type: Ready
  - some: thing
---
apiVersion: v1
kind: Node
metadata:
  name: node-ng1-bbb
  labels:
    node.deckhouse.io/group: ng1
  annotations:
    node.deckhouse.io/configuration-checksum: a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3
status:
  conditions:
  - some: thing
  - status: "True"
    type: Ready
`
		stateCloudProviderSecret = `
---
apiVersion: v1
kind: Secret
metadata:
  name: d8-node-manager-cloud-provider
  namespace: kube-system
data:
  zones: WyJub3ZhIl0= # ["nova"]
`
		configurationChecksums = `
---
apiVersion: v1
kind: Secret
metadata:
  name: configuration-checksums
  namespace: d8-cloud-instance-manager
data:
  ng1: YTY2NWE0NTkyMDQyMmY5ZDQxN2U0ODY3ZWZkYzRmYjhhMDRhMWYzZmZmMWZhMDdlOTk4ZTg2ZjdmN2EyN2FlMw== # sha256sum 123
  ng2: OGQyM2NmNmM4NmU4MzRhN2FhNmVkZWQ1NGMyNmNlMmJiMmU3NDkwMzUzOGM2MWJkZDVkMjE5Nzk5N2FiMmY3Mg== # sha256sum 321
`
	)

	f := HookExecutionConfigInit(`{}`, `{}`)
	f.RegisterCRD("deckhouse.io", "v1alpha1", "NodeGroup", false)
	f.RegisterCRD("machine.sapcloud.io", "v1alpha1", "MachineDeployment", true)
	f.RegisterCRD("machine.sapcloud.io", "v1alpha1", "Machine", true)

	Context("Empty cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(``))
			f.RunHook()
		})

		It("Hook must not fail", func() {
			Expect(f).To(ExecuteSuccessfully())
		})
	})

	Context("A NG1 and zones Secret", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateCloudNG1 + stateCloudProviderSecret))
			f.RunHook()
		})

		It("Min and max must be filled", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("NodeGroup", "ng1").Field("status").String()).To(MatchJSON(`{"extra":"thing","max":5,"min":1,"desired":0,"instances":0,"nodes":0,"ready":0,"upToDate": 0}`))
		})
	})

	Context("NGs MD, Machines, Nodes and zones Secret", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateCloudNG1 + stateCloudNG2 + stateMDs + stateMachines + stateNodes + stateCloudProviderSecret + configurationChecksums))
			f.RunHook()
		})

		It("Min, max, desired, instances, nodes, ready must be filled", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("NodeGroup", "ng1").Field("status").String()).To(MatchJSON(`{"extra":"thing","max":5,"min":1,"desired":2,"instances":2,"nodes":2,"ready":1,"upToDate": 2}`))
			Expect(f.KubernetesGlobalResource("NodeGroup", "ng2").Field("status").String()).To(MatchJSON(`{"max":9,"min":6,"desired":0,"instances":0,"nodes":0,"ready":0,"upToDate": 0}`))
		})
	})

	Context("NGs MD, Machines, Nodes and zones Secret", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateNG1 + stateNG2 + stateMDs + stateMachines + stateNodes + stateCloudProviderSecret + configurationChecksums))
			f.RunHook()
		})

		It("Nodes, ready must be filled", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("NodeGroup", "ng1").Field("status").String()).To(MatchJSON(`{"extra":"thing","nodes":2,"ready":1,"upToDate": 2}`))
			Expect(f.KubernetesGlobalResource("NodeGroup", "ng2").Field("status").String()).To(MatchJSON(`{"nodes":0,"ready":0,"upToDate": 0}`))
		})
	})
})

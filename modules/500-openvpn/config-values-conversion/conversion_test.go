package config_values_conversion

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/conversion"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("Module :: openvpn :: config values conversions :: v0.0.0", func() {
	f := SetupConverter(``)

	const migratedValues = `
inlet: ExternalIP
hostPort: 2222
`
	Context("config values already migrated", func() {
		BeforeEach(func() {
			f.ValuesSetFromYaml(".", migratedValues)
			f.Convert("v0.0.0")
		})

		It("should convert", func() {
			Expect(f.Error).ShouldNot(HaveOccurred())
			Expect(f.FinalVersion).Should(Equal("v1.0.0"))
			Expect(f.FinalValues.Get("storageClass").String()).Should(BeEmpty())
		})
	})

	const nonMigratedValues = `
inlet: ExternalIP
hostPort: 2222
storageClass: default
`
	Context("config values are non migrated", func() {
		BeforeEach(func() {
			f.ValuesSetFromYaml(".", nonMigratedValues)
			f.Convert("v0.0.0")
		})

		It("should convert", func() {
			Expect(f.Error).ShouldNot(HaveOccurred())
			Expect(f.FinalVersion).Should(Equal("v1.0.0"))
			Expect(f.FinalValues.Get("storageClass").String()).Should(BeEmpty())
		})
	})
})

// Test older values conversion to latest version.
var _ = Describe("Module :: openvpn :: config values conversions :: to latest", func() {
	f := SetupConverter(``)

	Context("v0.0.0", func() {
		const v0_0_0_Values = `
inlet: ExternalIP
hostPort: 2222
storageClass: default
`

		BeforeEach(func() {
			f.ValuesSetFromYaml(".", v0_0_0_Values)
			f.ConvertToLatest("v0.0.0")
		})

		It("should convert", func() {
			Expect(f.Error).ShouldNot(HaveOccurred())
		})
	})
})

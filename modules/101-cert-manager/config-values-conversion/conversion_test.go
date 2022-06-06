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

package config_values_conversion

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
)

func Test_conv_v0(t *testing.T) {
	g := NewWithT(t)
	conv := conversion.Registry().Get("cert-manager", "v0.0.0")

	g.Expect(conv).ShouldNot(BeNil())

	v0Vals := map[string]interface{}{
		"email": "admin@example.com",
		"nodeSelector": map[string]interface{}{
			"has-gpu": "true",
		},
		"cloudflareAPIToken": "APITOKEN",
		"cloudflareEmail":    "example@example.com",
	}

	newVer, newVals, err := conv.Convert("v0.0.0", v0Vals)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(newVer).Should(Equal("v1.0.0"))
	g.Expect(newVals).Should(HaveKey("cloudflare"))
	cloudflare, ok := newVals["cloudflare"].(map[string]interface{})
	g.Expect(ok).Should(BeTrue(), "cloudflare section should be a map")
	g.Expect(cloudflare).Should(HaveKey("APIToken"))
	g.Expect(cloudflare["APIToken"]).Should(Equal(v0Vals["cloudflareAPIToken"]))
	g.Expect(cloudflare).Should(HaveKey("Email"))
	g.Expect(cloudflare["Email"]).Should(Equal(v0Vals["cloudflareEmail"]))
}

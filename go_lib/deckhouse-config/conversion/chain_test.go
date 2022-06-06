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

package conversion

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_conv_to_latest(t *testing.T) {
	g := NewWithT(t)

	const modName = "test-mod"
	RegisterFunc(modName, "v0.0.0", func(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error) {
		configValues["param2"] = "val2"
		return "v1.0.0", configValues, nil
	})
	RegisterFunc(modName, "v1.0.0", func(configVersion string, configValues map[string]interface{}) (string, map[string]interface{}, error) {
		return configVersion, configValues, nil
	})

	v0Vals := map[string]interface{}{
		"param1": "val1",
	}
	newVer, newVals, err := ConvertToLatest(modName, "v0.0.0", v0Vals)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(newVer).Should(Equal("v1.0.0"))
	g.Expect(newVals).Should(HaveKey("param1"), "should keep old params")
	g.Expect(newVals).Should(HaveKey("param2"), "should add new param")
}

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

package matrix

import (
	"os"
	"strings"
	"testing"

	"github.com/deckhouse/deckhouse/testing/matrix/linter"
	"github.com/deckhouse/deckhouse/testing/matrix/linter/rules/modules"
	"github.com/stretchr/testify/require"
)

func TestMatrix(t *testing.T) {
	discoveredModules, err := modules.GetDeckhouseModulesWithValuesMatrixTests()
	require.NoError(t, err)

	// Use environment variable to focus on specific module, e.g. D8_TEST_MATRIX_FOCUS=user-authn,user-authz
	focus := os.Getenv("FOCUS")

	focusNames := make(map[string]struct{})
	if focus != "" {
		parts := strings.Split(focus, ",")
		for _, part := range parts {
			focusNames[part] = struct{}{}
		}
	}

	changeSymlinks(t)
	for _, module := range discoveredModules {
		_, ok := focusNames[module.Name]
		if len(focusNames) == 0 || ok {
			t.Run(module.Name, func(t *testing.T) {
				require.NoError(t, linter.Run("", module))
			})
		}
	}
	restoreSymlinks(t)
}

// changeSymlinks changes symlinks in module dir to proper place when modules in ee/fe not copied to main modules directory
func changeSymlink(t *testing.T, symlinkPath string, newDestination string) {
	_, err := os.Lstat(symlinkPath)
	require.NoError(t, err)

	err = os.Remove(symlinkPath)
	require.NoError(t, err)

	err = os.Symlink(newDestination, symlinkPath)
	require.NoError(t, err)
}

func changeSymlinks(t *testing.T) {
	changeSymlink(t, "/deckhouse/ee/modules/030-cloud-provider-openstack/candi", "/deckhouse/ee/candi/cloud-providers/openstack/")
	changeSymlink(t, "/deckhouse/ee/modules/030-cloud-provider-vsphere/candi", "/deckhouse/ee/candi/cloud-providers/vsphere/")

	_, err := os.Lstat("/deckhouse/modules/040-node-manager/images_tags.json")
	require.NoError(t, err)
	err = os.Remove("/deckhouse/modules/040-node-manager/images_tags.json")
	require.NoError(t, err)
	err = os.Symlink("/deckhouse/ee/modules/030-cloud-provider-openstack/cloud-instance-manager/", "/deckhouse/modules/040-node-manager/cloud-providers/openstack")
	require.NoError(t, err)
	err = os.Symlink("/deckhouse/ee/modules/030-cloud-provider-vsphere/cloud-instance-manager/", "/deckhouse/modules/040-node-manager/cloud-providers/vsphere")
	require.NoError(t, err)
}

// restoreSymlinks restores symlinks in module dir to original place
func restoreSymlinks(t *testing.T) {
	changeSymlink(t, "/deckhouse/ee/modules/030-cloud-provider-openstack/candi", "/deckhouse/candi/cloud-providers/openstack/")
	changeSymlink(t, "/deckhouse/ee/modules/030-cloud-provider-vsphere/candi", "/deckhouse/candi/cloud-providers/vsphere/")
	err := os.Symlink("../images_tags.json", "/deckhouse/modules/040-node-manager/images_tags.json")
	require.NoError(t, err)
	err = os.Remove("/deckhouse/modules/040-node-manager/cloud-providers/openstack")
	require.NoError(t, err)
	err = os.Remove("/deckhouse/modules/040-node-manager/cloud-providers/vsphere")
	require.NoError(t, err)
}

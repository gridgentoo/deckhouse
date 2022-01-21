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

package change_host_address

import (
	"github.com/flant/shell-operator/pkg/kube/object_patch"
)

type podClient interface {
	Delete(name string)
	AnnotateHost(name, host string)
}

type podClientImpl struct {
	annoKey   string
	namespace string

	patcher *object_patch.PatchCollector
}

func (pc *podClientImpl) AnnotateHost(name string, host string) {
	patch := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				pc.annoKey: host,
			},
		},
	}
	pc.patcher.MergePatch(patch, "v1", "Pod", pc.namespace, name)
}

func (pc *podClientImpl) Delete(name string) {
	pc.patcher.Delete("v1", "Pod", pc.namespace, name)
}

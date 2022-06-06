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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeckhouseConfig is a configuration for module or for global config values.
type DeckhouseConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DeckhouseConfigSpec `json:"spec"`

	Status DeckhouseConfigStatus `json:"status,omitempty"`
}

type DeckhouseConfigSpec struct {
	ConfigVersion string                 `json:"configVersion,omitempty"`
	ConfigValues  map[string]interface{} `json:"configValues,omitempty"`
	Enabled       *bool                  `json:"enabled,omitempty"`
}

type DeckhouseConfigStatus struct {
	Enabled string `json:"enabled"`
	Status  string `json:"status"`
}

type deckhouseConfigKind struct{}

func (in *DeckhouseConfigStatus) GetObjectKind() schema.ObjectKind {
	return &deckhouseConfigKind{}
}

func (f *deckhouseConfigKind) SetGroupVersionKind(_ schema.GroupVersionKind) {}
func (f *deckhouseConfigKind) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: "deckhouse.io", Version: "v1", Kind: "DeckhouseConfig"}
}

func GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "deckhouse.io",
		Version:  "v1",
		Resource: "deckhouseconfigs",
	}
}

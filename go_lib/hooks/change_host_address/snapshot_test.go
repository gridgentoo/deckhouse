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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func Test_getAddress(t *testing.T) {

	unstructuredPod := genPod(podArgs{
		Name:      "podname",
		IP:        "1.2.3.4",
		InitialIP: "4.3.2.1",
	})

	expected := podHost{
		Name:      "podname",
		IP:        "1.2.3.4",
		InitialIP: "4.3.2.1",
	}

	got, err := getAddress(unstructuredPod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equalf(t, expected, got.(podHost), "getAddress(%v)", unstructuredPod)

}

type podArgs struct {
	Name      string
	IP        string
	InitialIP string
}

func genPod(args podArgs) *unstructured.Unstructured {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: args.Name,
		},
		Status: v1.PodStatus{HostIP: args.IP},
	}

	if args.InitialIP != "" {
		pod.ObjectMeta.SetAnnotations(map[string]string{
			"node.deckhouse.io/initial-host-ip": args.InitialIP,
		})
	}

	manifest, err := json.Marshal(pod)
	if err != nil {
		panic(err)
	}

	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	if _, _, err := decUnstructured.Decode(manifest, nil, obj); err != nil {
		panic(err)
	}

	return obj
}

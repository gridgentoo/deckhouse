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
	"fmt"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type podHost struct {
	Name      string
	IP        string
	InitialIP string
}

func getAddress(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	pod := &v1.Pod{}
	err := sdk.FromUnstructured(obj, pod)
	if err != nil {
		return nil, fmt.Errorf("cannot convert pod: %v", err)
	}

	return podHost{
		Name:      pod.Name,
		IP:        pod.Status.HostIP,
		InitialIP: pod.Annotations[initialHostAddressAnnotation],
	}, nil
}

func parsePods(snaps []go_hook.FilterResult) []podHost {
	ps := make([]podHost, len(snaps))
	for i, s := range snaps {
		ps[i] = s.(podHost)
	}
	return ps
}

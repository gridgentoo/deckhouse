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
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const initialHostAddressAnnotation = "node.deckhouse.io/initial-host-ip"

func RegisterHook(appName, namespace string) bool {
	return sdk.RegisterFunc(&go_hook.HookConfig{
		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:       "pod",
				ApiVersion: "v1",
				Kind:       "Pod",
				NamespaceSelector: &types.NamespaceSelector{
					NameSelector: &types.NameSelector{
						MatchNames: []string{namespace},
					},
				},
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Key:      "app",
							Operator: "In",
							Values:   []string{appName},
						},
					},
				},
				FilterFunc: getAddress,
			},
		},
	}, wrapChangeAddressHandler(namespace))
}

func wrapChangeAddressHandler(namespace string) func(input *go_hook.HookInput) error {
	return func(input *go_hook.HookInput) error {
		pods := parsePods(input.Snapshots["pod"])
		client := &podClientImpl{
			annoKey:   initialHostAddressAnnotation,
			namespace: namespace,
			patcher:   input.PatchCollector,
		}

		changeHostAddress(client, pods)

		return nil
	}
}

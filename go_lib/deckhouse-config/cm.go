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

package deckhouse_config

import (
	"context"

	v1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/deckhouse/deckhouse/go_lib/dependency/k8s"
)

const GeneratedConfigMapName = "deckhouse-generated-config-do-not-edit"
const DeckhouseConfigMapName = "deckhouse"
const DeckhouseNS = "d8-system"

// HasGeneratedConfigMap returns true if generated config is present in cluster.
func HasGeneratedConfigMap(klient k8s.Client) (bool, error) {
	cm, err := GetGeneratedConfigMap(klient)
	if k8errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return cm != nil, nil
}

// GetGeneratedConfigMap returns generated ConfigMap with config values.
func GetGeneratedConfigMap(klient k8s.Client) (*v1.ConfigMap, error) {
	return klient.CoreV1().ConfigMaps(DeckhouseNS).Get(context.TODO(), GeneratedConfigMapName, metav1.GetOptions{})
}

// GetDeckhouseConfigMap returns default ConfigMap with config values (ConfigMap/deckhouse).
func GetDeckhouseConfigMap(klient k8s.Client) (*v1.ConfigMap, error) {
	return klient.CoreV1().ConfigMaps(DeckhouseNS).Get(context.TODO(), DeckhouseConfigMapName, metav1.GetOptions{})
}

func GeneratedConfigMap(data map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GeneratedConfigMapName,
			Namespace: DeckhouseNS,
			Labels: map[string]string{
				"owner": "deckhouse",
			},
		},
		Data: data,
	}
}

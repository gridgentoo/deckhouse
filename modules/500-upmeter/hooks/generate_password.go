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

package hooks

import (
	"fmt"
	"strings"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/deckhouse/deckhouse/go_lib/pwgen"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: "/modules/upmeter/generate_password",
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       authSecretBinding,
			ApiVersion: "v1",
			Kind:       "Secret",
			NameSelector: &types.NameSelector{
				MatchNames: authSecretNames,
			},
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{upmeterNS},
				},
			},
			// Synchronization is redundant because of OnBeforeHelm.
			ExecuteHookOnSynchronization: go_hook.Bool(false),
			ExecuteHookOnEvents:          go_hook.Bool(false),
			FilterFunc:                   filterAuthSecret,
		},
	},

	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
}, restoreOrGeneratePassword)

const (
	upmeterNS         = "d8-upmeter"
	authSecretField   = "auth"
	authSecretBinding = "auth-secrets"
	statusSecretName  = "basic-auth-status"
	webuiSecretName   = "basic-auth-webui"

	externalAuthValuesTmpl     = "upmeter.auth.%s.externalAuthentication"
	passwordValuesTmpl         = "upmeter.auth.%s.password"
	passwordInternalValuesTmpl = "upmeter.internal.auth.%s.password"
)

var authSecretNames = []string{statusSecretName, webuiSecretName}
var upmeterApps = map[string]string{
	statusSecretName: "status",
	webuiSecretName:  "webui",
}

type storedPassword struct {
	SecretName string `json:"name"`
	Password   string `json:"password,omitempty"`
}

func filterAuthSecret(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	secret := &v1.Secret{}
	err := sdk.FromUnstructured(obj, secret)
	if err != nil {
		return nil, fmt.Errorf("cannot convert secret to struct: %v", err)
	}

	password := string(secret.Data[authSecretField])
	if password != "" {
		password = strings.TrimPrefix(password, "auth:{PLAIN}")
	}

	return storedPassword{
		SecretName: secret.GetName(),
		Password:   password,
	}, nil
}

// restoreOrGeneratePassword restores passwords from config values or secrets.
// If there are no passwords, it generates new.
func restoreOrGeneratePassword(input *go_hook.HookInput) error {
	storedPasswords := map[string]string{}
	for _, snap := range input.Snapshots[authSecretBinding] {
		storedPassword := snap.(storedPassword)
		storedPasswords[storedPassword.SecretName] = storedPassword.Password
	}

	for secretName, appName := range upmeterApps {
		externalAuthValuesPath := fmt.Sprintf(externalAuthValuesTmpl, appName)
		passwordValuesPath := fmt.Sprintf(passwordValuesTmpl, appName)
		passwordInternalValuesPath := fmt.Sprintf(passwordInternalValuesTmpl, appName)

		// Clear password from internal values if an external authentication is enabled.
		if input.Values.Exists(externalAuthValuesPath) {
			input.Values.Remove(passwordInternalValuesPath)
			continue
		}

		// Try to set password from config values.
		password, ok := input.ConfigValues.GetOk(passwordValuesPath)
		if ok {
			input.Values.Set(passwordInternalValuesPath, password.String())
			continue
		}

		// Try to set password from the stored password.
		if storedPassword, has := storedPasswords[secretName]; has {
			input.Values.Set(passwordInternalValuesPath, storedPassword)
			continue
		}

		// Password is already in values. It should not happen, just a precaution.
		_, ok = input.Values.GetOk(passwordInternalValuesPath)
		if ok {
			continue
		}

		// No password found for the app, generate new one.
		newPasswd := pwgen.AlphaNum(20)
		input.Values.Set(passwordInternalValuesPath, newPasswd)
	}

	return nil
}

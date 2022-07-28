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

package generate_password

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

const (
	secretBindingName          = "password_secret"
	defaultBasicAuthPlainField = "auth"
	defaultBeforeHelmOrder     = 10
)

func NewBasicAuthPlainHook(moduleValuesPath string, ns string, secretName string) *Hook {
	return &Hook{
		Secret: Secret{
			Namespace:           ns,
			Name:                secretName,
			BasicAuthPlainField: defaultBasicAuthPlainField,
		},
		Values: Values{
			ModuleKey: moduleValuesPath,
		},
		BeforeHelmOrder: defaultBeforeHelmOrder,
	}
}

// RegisterHook returns func to register common hook that generates
// and stores a password in the Secret.
func RegisterHook(hook *Hook) bool {
	return sdk.RegisterFunc(&go_hook.HookConfig{
		Queue: fmt.Sprintf("/modules/%s/generate_password", hook.Values.ModuleKey),
		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:       secretBindingName,
				ApiVersion: "v1",
				Kind:       "Secret",
				NameSelector: &types.NameSelector{
					MatchNames: []string{hook.Secret.Name},
				},
				NamespaceSelector: &types.NamespaceSelector{
					NameSelector: &types.NameSelector{
						MatchNames: []string{hook.Secret.Namespace},
					},
				},
				// Synchronization is redundant because of OnBeforeHelm.
				ExecuteHookOnSynchronization: go_hook.Bool(false),
				ExecuteHookOnEvents:          go_hook.Bool(false),
				FilterFunc:                   hook.Filter,
			},
		},
		OnBeforeHelm: &go_hook.OrderedConfig{Order: float64(hook.BeforeHelmOrder)},
	}, hook.Handle)
}

type Hook struct {
	Secret          Secret
	Values          Values
	BeforeHelmOrder int

	PasswordGeneratorFunc func(input *go_hook.HookInput) string
}

type Secret struct {
	Namespace string
	Name      string

	RawPasswordField    string
	BasicAuthPlainField string

	FilterFunc func(secret *v1.Secret) (go_hook.FilterResult, error)
}

type Values struct {
	ModuleKey           string
	PasswordKey         string
	PasswordInternalKey string
	ExternalAuthKey     string
}

// Filter extracts password from the Secret. Password can be stored as a raw string or as
// a basic auth plain format (user:{PLAIN}password). Custom FilterFunc is called for custom
// password extraction.
func (h *Hook) Filter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	secret := &v1.Secret{}
	err := sdk.FromUnstructured(obj, secret)
	if err != nil {
		return nil, fmt.Errorf("cannot convert secret to struct: %v", err)
	}

	if h.Secret.FilterFunc != nil {
		return h.Secret.FilterFunc(secret)
	}

	if h.Secret.RawPasswordField != "" {
		return string(secret.Data[h.Secret.RawPasswordField]), nil
	}

	// Default is to store password as basic auth plain format.
	field := h.Secret.BasicAuthPlainField
	if field == "" {
		field = defaultBasicAuthPlainField
	}

	basicAuth := string(secret.Data[field])
	if !strings.Contains(basicAuth, "{PLAIN}") {
		return nil, fmt.Errorf("field '%s' in Secret/%s is not a basic auth plain password", field, secret.GetName())
	}

	parts := strings.SplitN(basicAuth, "{PLAIN}", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("field '%s' in Secret/%s is not a basic auth plain password", field, secret.GetName())
	}
	return strings.TrimSpace(parts[1]), nil
}

// Handle restores password from the configuration or from the Secret and
// puts it to internal values.
// It generates new password if no password found in the configuration and the
// Secret or no externalAuthentication defined.
func (h *Hook) Handle(input *go_hook.HookInput) error {
	var externalAuthKey = h.ExternalAuthKey()
	var passwordKey = h.PasswordKey()
	var passwordInternalKey = h.PasswordInternalKey()

	// Clear password from internal values if an external authentication is enabled.
	if input.Values.Exists(externalAuthKey) {
		input.Values.Remove(passwordInternalKey)
		return nil
	}

	// Try to set password from config values.
	password, ok := input.ConfigValues.GetOk(passwordKey)
	if ok {
		input.Values.Set(passwordInternalKey, password.String())
		return nil
	}

	// Try to set password from the Secret.
	snap := input.Snapshots[secretBindingName]
	if len(snap) > 0 {
		storedPassword := snap[0].(string)
		input.Values.Set(passwordInternalKey, storedPassword)
		return nil
	}

	// Return if auth key is already in values.
	_, ok = input.Values.GetOk(passwordInternalKey)
	if ok {
		return nil
	}

	// No password found, generate new one.
	newPassword := ""
	if h.PasswordGeneratorFunc != nil {
		newPassword = h.PasswordGeneratorFunc(input)
	} else {
		newPassword = DefaultGenerator()
	}

	input.Values.Set(passwordInternalKey, newPassword)
	return nil
}

const (
	externalAuthKeyTmpl     = "%s.auth.externalAuthentication"
	passwordKeyTmpl         = "%s.auth.password"
	passwordInternalKeyTmpl = "%s.internal.auth.password"
)

func (h *Hook) ExternalAuthKey() string {
	key := h.Values.ExternalAuthKey
	if key != "" {
		return key
	}
	return fmt.Sprintf(externalAuthKeyTmpl, h.Values.ModuleKey)
}

func (h *Hook) PasswordKey() string {
	key := h.Values.PasswordKey
	if key != "" {
		return key
	}
	return fmt.Sprintf(passwordKeyTmpl, h.Values.ModuleKey)
}

func (h *Hook) PasswordInternalKey() string {
	key := h.Values.PasswordInternalKey
	if key != "" {
		return key
	}
	return fmt.Sprintf(passwordInternalKeyTmpl, h.Values.ModuleKey)
}

func DefaultGenerator() string {
	return pwgen.AlphaNum(20)
}

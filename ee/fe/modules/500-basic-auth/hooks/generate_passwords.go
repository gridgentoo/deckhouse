/*
Copyright 2021 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"fmt"
	"regexp"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/deckhouse/deckhouse/go_lib/pwgen"
)

// Set locations from config values or a default one with generated password.

// Subscribe to Secret/htpasswd from templates.
const secretNS = "kube-basic-auth"
const secretName = "htpasswd"
const secretBinding = "htpasswd_secret"
const locationsKey = "basicAuth.locations"
const locationsInternalKey = "basicAuth.internal.locations"
const generatedPasswdLength = 20

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       secretBinding,
			ApiVersion: "v1",
			Kind:       "Secret",
			NameSelector: &types.NameSelector{
				MatchNames: []string{secretName},
			},
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{secretNS},
				},
			},
			// Synchronization is redundant because of OnBeforeHelm.
			ExecuteHookOnSynchronization: go_hook.Bool(false),
			ExecuteHookOnEvents:          go_hook.Bool(false),
			FilterFunc:                   filterHtpasswdSecret,
		},
	},

	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
}, generatePassword)

func filterHtpasswdSecret(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	secret := &v1.Secret{}
	err := sdk.FromUnstructured(obj, secret)
	if err != nil {
		return nil, fmt.Errorf("cannot convert secret to struct: %v", err)
	}

	return string(secret.Data["htpasswd"]), nil
}

const defaultUserName = `admin`

func defaultLocationValues(passwd string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"users": map[string]interface{}{
				defaultUserName: passwd,
			},
			"location": "/",
		},
	}
}

// Regex to get password from basic auth plain format.
// Format is: username:{PLAIN}password
// htpasswd field contains several passwords, so use m flag for multi-line mode.
var defaultPasswordRe = regexp.MustCompile(`(?m)^\s*` + defaultUserName + `:{PLAIN}(\S+)$`)

func generatePassword(input *go_hook.HookInput) error {
	// Set values from user controlled configuration.
	userLocations, ok := input.ConfigValues.GetOk(locationsKey)
	if ok {
		input.Values.Set(locationsInternalKey, userLocations.Value())
		return nil
	}

	_, ok = input.Values.GetOk(locationsInternalKey)
	if ok {
		// No config values, but internal values are set. It's OK, just return.
		return nil
	}

	// No values, no secret. Module is enabled for the first time, so
	// generate a new password and prepare a default location.
	if len(input.Snapshots[secretBinding]) == 0 {
		locations := defaultLocationValues(pwgen.AlphaNum(generatedPasswdLength))
		input.Values.Set(locationsInternalKey, locations)
		return nil
	}

	// No values, but Secret is present. This can occur when module is enabled
	// and Deckhouse is restarted later. Restore generated password from the Secret
	// assuming it is in the first line of htpasswd field.
	// NOTE: This algorithm is coupled with the field name in secret.yaml and "users" template in _helpers.tpl.
	htpasswdField := input.Snapshots[secretBinding][0].(string)
	matches := defaultPasswordRe.FindStringSubmatch(htpasswdField)
	// matches[0] is a full string
	// matches[1] is a password
	if len(matches) != 2 || len(matches[1]) != generatedPasswdLength {
		return fmt.Errorf("expect secret/%s contains generated credentials in basic auth plain format (remove Secret to generate new credentials)", secretName)
	}

	locations := defaultLocationValues(matches[1])
	input.Values.Set(locationsInternalKey, locations)
	return nil
}

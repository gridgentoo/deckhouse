package internal

import (
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"

	d8config_v1 "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/v1"
	"github.com/deckhouse/deckhouse/go_lib/set"
)

func SetFromArrayValue(input *go_hook.HookInput, path string) (set.Set, error) {
	value, ok := input.Values.GetOk(path)
	if !ok {
		return nil, fmt.Errorf("%s value is required", path)
	}
	list := value.Array()
	if len(list) == 0 {
		return nil, fmt.Errorf("%s value should not be empty", path)
	}
	return set.NewFromValues(input.Values, path), nil
}

func KnownConfigsFromSnapshot(snapshot []go_hook.FilterResult, possibleNames set.Set) []*d8config_v1.DeckhouseConfig {
	configs := make([]*d8config_v1.DeckhouseConfig, 0)
	for _, item := range snapshot {
		cfg := item.(*d8config_v1.DeckhouseConfig)
		// Ignore unknown names.
		if !possibleNames.Has(cfg.GetName()) {
			continue
		}
		configs = append(configs, cfg)
	}
	return configs
}

func ConfigsFromSnapshot(snapshot []go_hook.FilterResult) []*d8config_v1.DeckhouseConfig {
	configs := make([]*d8config_v1.DeckhouseConfig, 0)
	for _, item := range snapshot {
		cfg := item.(*d8config_v1.DeckhouseConfig)
		configs = append(configs, cfg)
	}
	return configs
}

// MergeEnabled merges enabled flags. Enabled flag can be nil.
//
// If all flags are nil, then false is returned — module is disabled by default.
// Note: copy-paste from AddonOperator.ModuleManager
func MergeEnabled(enabledFlags ...*bool) bool {
	result := false
	for _, enabled := range enabledFlags {
		if enabled == nil {
			continue
		} else {
			result = *enabled
		}
	}

	return result
}

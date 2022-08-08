package hooks

import (
	"os"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnStartup: &go_hook.OrderedConfig{Order: 5},
}, webhookSettings)

func webhookSettings(input *go_hook.HookInput) error {
	input.Values.Set("deckhouse.internal.configMapName", os.Getenv("ADDON_OPERATOR_CONFIG_MAP"))
	return nil
}

package internal

import (
	"pb_launcher/internal/hooks"

	"go.uber.org/fx"
)

var Module = fx.Module("hooks",
	fx.Invoke(hooks.RegisterAdminExistsRoute),
	fx.Invoke(hooks.RegisterServiceLogsRoute),
	fx.Invoke(hooks.AddServiceHooks),
	fx.Invoke(hooks.AddComandHooks),
)

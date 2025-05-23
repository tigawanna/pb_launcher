package hooks

import (
	"pb_launcher/collections"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func AddComandHooks(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest(collections.ServicesComands).
		BindFunc(func(e *core.RecordRequestEvent) error {
			e.Record.Set("status", "pending")
			e.Record.Set("error_message", nil)
			e.Record.Set("executed", nil)
			return e.Next()
		})
}

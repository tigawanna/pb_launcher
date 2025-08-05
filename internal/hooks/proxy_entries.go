package hooks

import (
	"pb_launcher/collections"
	"pb_launcher/internal/proxy/domain"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func AddProxyEntriesHooks(app *pocketbase.PocketBase,
	discovery *domain.ProxyEntryDiscovery,
) {
	app.OnRecordAfterUpdateSuccess(collections.ProxyEntries).
		BindFunc(func(e *core.RecordEvent) error {
			if err := e.Next(); err != nil {
				return err
			}
			discovery.InvalidateProxyEntriesCacheByID(e.Record.Id)
			return nil
		})
}

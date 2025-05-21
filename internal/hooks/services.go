package hooks

import (
	"errors"
	"pb_launcher/collections"
	"slices"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func AddServiceHooks(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest(collections.Services).
		BindFunc(func(e *core.RecordRequestEvent) error {
			if e.Auth == nil {
				return errors.New("unauthorized: no auth record found")
			}

			email := e.Auth.GetString("email")
			if email == "" {
				return errors.New("unauthorized: email missing in auth record")
			}

			restart_policy := e.Record.GetString("restart_policy")
			if !slices.Contains([]string{}, restart_policy) {
				restart_policy = "no"
			}

			e.Record.Set("boot_completed", "no")
			e.Record.Set("restart_policy", restart_policy)
			e.Record.Set("boot_user_email", email)
			e.Record.Set("boot_user_password", core.GenerateDefaultRandomId())
			e.Record.Set("status", "idle")

			return e.Next()
		})

	app.OnRecordUpdateRequest(collections.Services).BindFunc(func(e *core.RecordRequestEvent) error {
		updatedName := e.Record.GetString("name")
		updatedPolicy := e.Record.Get("restart_policy")
		isDeleted := e.Record.Get("deleted")

		currentRecord, err := e.App.FindRecordById(e.Collection, e.Record.GetString("id"))
		if err != nil {
			return err
		}

		currentRecord.Set("name", updatedName)
		currentRecord.Set("restart_policy", updatedPolicy)
		currentRecord.Set("deleted", isDeleted)

		e.Record = currentRecord
		return e.Next()
	})

	app.OnRecordAfterCreateSuccess(collections.Services).BindFunc(func(e *core.RecordEvent) error {
		comandCollection, err := e.App.FindCachedCollectionByNameOrId(collections.ServicesComands)
		if err != nil {
			return err
		}
		record := core.NewRecord(comandCollection)

		record.Set("service", e.Record.Get("id"))
		record.Set("action", "start")
		record.Set("status", "pending")
		record.Set("error_message", "")
		record.Set("executed", nil)

		if err := e.App.Save(record); err != nil {
			return err
		}
		return e.Next()
	})
}

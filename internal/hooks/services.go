package hooks

import (
	"errors"
	"pb_launcher/collections"
	certmanager "pb_launcher/internal/certmanager/domain"
	"pb_launcher/internal/certmanager/domain/repositories"
	"pb_launcher/internal/proxy/domain"
	"slices"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func AddServiceHooks(app *pocketbase.PocketBase,
	serviceDiscovery *domain.ServiceDiscovery,
	domainDiscovery *domain.DomainServiceDiscovery,
	planner *certmanager.CertRequestPlannerUsecase,
	repository repositories.CertRequestRepository,
) {
	app.OnRecordCreateRequest(collections.Services).
		BindFunc(func(e *core.RecordRequestEvent) error {
			if e.Auth == nil {
				return errors.New("unauthorized: no auth record found")
			}

			restart_policy := e.Record.GetString("restart_policy")
			if !slices.Contains([]string{"no", "on-failure"}, restart_policy) {
				restart_policy = "no"
			}

			e.Record.Set("boot_completed", "no")
			e.Record.Set("restart_policy", restart_policy)

			e.Record.Set("status", "idle")

			return e.Next()
		})

	app.OnRecordUpdateRequest(collections.Services).BindFunc(func(e *core.RecordRequestEvent) error {
		updatedName := e.Record.GetString("name")
		updatedPolicy := e.Record.Get("restart_policy")
		deleted := e.Record.GetDateTime("deleted")

		currentRecord, err := e.App.FindRecordById(e.Collection, e.Record.GetString("id"))
		if err != nil {
			return err
		}

		currentRecord.Set("name", updatedName)
		currentRecord.Set("restart_policy", updatedPolicy)
		currentRecord.Set("deleted", deleted)

		e.Record = currentRecord
		if err := e.Next(); err != nil {
			return err
		}
		if !deleted.IsZero() {
			comandCollection, err := e.App.FindCachedCollectionByNameOrId(collections.ServicesComands)
			if err != nil {
				return err
			}
			record := core.NewRecord(comandCollection)

			record.Set("service", e.Record.Id)
			record.Set("action", "stop")
			record.Set("status", "pending")
			record.Set("error_message", "")
			record.Set("executed", nil)

			if err := e.App.Save(record); err != nil {
				return err
			}
		}
		return nil
	})

	app.OnRecordAfterCreateSuccess(collections.Services).BindFunc(func(e *core.RecordEvent) error {
		comandCollection, err := e.App.FindCachedCollectionByNameOrId(collections.ServicesComands)
		if err != nil {
			return err
		}
		record := core.NewRecord(comandCollection)

		record.Set("service", e.Record.Id)
		record.Set("action", "start")
		record.Set("status", "pending")
		record.Set("error_message", "")
		record.Set("executed", nil)

		if err := e.App.Save(record); err != nil {
			return err
		}
		return e.Next()
	})

	app.OnRecordAfterUpdateSuccess(collections.Services).
		BindFunc(func(e *core.RecordEvent) error {
			if err := e.Next(); err != nil {
				return err
			}
			serviceDiscovery.InvalidateServiceCacheByID(e.Record.Id)
			return nil
		})

	app.OnRecordAfterCreateSuccess(collections.ServicesDomains).BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		domain := e.Record.GetString("domain")
		if e.Record.GetString("use_https") == "yes" {
			return planner.PostSSLDomainRequest(e.Context, domain, false)
		}
		return nil
	})

	app.OnRecordAfterUpdateSuccess(collections.ServicesDomains).
		BindFunc(func(e *core.RecordEvent) error {
			if err := e.Next(); err != nil {
				return err
			}
			domain := e.Record.GetString("domain")
			domainDiscovery.InvalidateDomain(domain)
			if e.Record.GetString("use_https") == "yes" {
				return planner.PostSSLDomainRequest(e.Context, domain, false)
			}
			return nil
		})

	app.OnRecordAfterDeleteSuccess(collections.ServicesDomains).
		BindFunc(func(e *core.RecordEvent) error {
			if err := e.Next(); err != nil {
				return err
			}
			domain := e.Record.GetString("domain")
			domainDiscovery.InvalidateDomain(domain)
			return repository.DeletePendingByDomain(e.Context, domain)
		})

}

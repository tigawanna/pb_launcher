package hooks

import (
	"pb_launcher/collections"
	"pb_launcher/configs"
	"pb_launcher/internal/certificates/tlscommon"
	certmanager "pb_launcher/internal/certmanager/domain"
	"pb_launcher/internal/certmanager/domain/repositories"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func AddServiceDomainsHooks(
	app *pocketbase.PocketBase,
	repository repositories.CertRequestRepository,
	planner *certmanager.CertRequestPlannerUsecase,
	store tlscommon.Store,
	cnf configs.Config,
) {
	app.OnRecordsListRequest(collections.ServicesDomains).BindFunc(
		func(e *core.RecordsListRequestEvent) error {
			baseCollecion, err := e.App.FindCollectionByNameOrId(collections.ServicesDomains)
			if err != nil {
				return nil
			}
			baseCollecion.Fields.Add(&core.TextField{Name: "x_cert_request_state"})
			baseCollecion.Fields.Add(&core.BoolField{Name: "x_reached_max_attempt"})
			baseCollecion.Fields.Add(&core.TextField{Name: "x_failed_error_message"})
			baseCollecion.Fields.Add(&core.BoolField{Name: "x_has_valid_ssl_cert"})

			for idx, record := range e.Records {
				if record.GetString("use_https") != "yes" {
					continue
				}
				domain := record.GetString("domain")
				last, err := repository.LastByDomain(e.Request.Context(), domain)
				if err != nil || last == nil {
					continue
				}
				newRecord := core.NewRecord(baseCollecion)
				for _, field := range record.Collection().Fields {
					fieldName := field.GetName()
					fieldValue := record.Get(fieldName)
					newRecord.Set(fieldName, fieldValue)
				}

				newRecord.Set("x_cert_request_state", string(last.Status))
				newRecord.Set("x_reached_max_attempt", last.Attempt >= cnf.GetMaxDomainCertAttempts())
				newRecord.Set("x_failed_error_message", last.Message)

				cert, err := store.Resolve(domain)
				if err != nil || cert == nil {
					newRecord.Set("x_has_valid_ssl_cert", false)
				} else {
					newRecord.Set("x_has_valid_ssl_cert", cert.Ttl > 0)
				}
				e.Records[idx] = newRecord
			}
			return e.Next()
		})
	app.OnRecordCreateRequest(collections.CertRequests).BindFunc(func(e *core.RecordRequestEvent) error {
		return planner.PostSSLDomainRequest(e.Request.Context(),
			e.Record.GetString("domain"),
			false,
		)
	})
}

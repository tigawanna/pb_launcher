package hooks

import (
	"net/http"
	"pb_launcher/helpers/logstore"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterServiceLogsRoute(app *pocketbase.PocketBase, lstore *logstore.ServiceLogDB) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/x-api/service/logs/{service_id}", func(re *core.RequestEvent) error {
			if re.Auth == nil {
				return re.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			serviceID := re.Request.PathValue("service_id")
			logs, err := lstore.GetLogsByService(serviceID)
			if err != nil {
				return re.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			return re.JSON(http.StatusOK, logs)
		})

		return e.Next()
	})
}

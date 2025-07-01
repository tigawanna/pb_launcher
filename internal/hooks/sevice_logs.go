package hooks

import (
	"net/http"
	"pb_launcher/helpers/logstore"
	"strconv"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterServiceLogsRoute(app *pocketbase.PocketBase, store *logstore.ServiceLogDB) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/x-api/service/logs/{service_id}", handleGetServiceLogs(store))
		e.Router.GET("/x-api/service/logs/{service_id}/{limit}", handleGetServiceLogs(store))
		return e.Next()
	})
}

func handleGetServiceLogs(store *logstore.ServiceLogDB) func(re *core.RequestEvent) error {
	return func(re *core.RequestEvent) error {
		if re.Auth == nil {
			return re.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		serviceID := re.Request.PathValue("service_id")
		limitStr := re.Request.PathValue("limit")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limitStr == "" {
			limit = -1
		}

		logs, err := store.GetLogsByService(serviceID, limit)
		if err != nil {
			return re.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return re.JSON(http.StatusOK, logs)
	}
}

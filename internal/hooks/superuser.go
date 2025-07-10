package hooks

import (
	"errors"
	"net/http"

	launcher "pb_launcher/internal/launcher/domain"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterUpsertServiceSuperuserRoute(app *pocketbase.PocketBase, launcherManager *launcher.LauncherManager) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/x-api/upsert_superuser/{service_id}",
			func(re *core.RequestEvent) error {
				email := re.Auth.GetString("email")
				if email == "" {
					return errors.New("unauthorized: email missing in auth record")
				}
				serviceID := re.Request.PathValue("service_id")
				password := core.GenerateDefaultRandomId()

				err := launcherManager.UpsertSuperuser(re.Request.Context(),
					serviceID, email, password)

				if err != nil {
					return re.InternalServerError("failed to upsert superuser", nil)
				}

				return re.JSON(http.StatusOK, map[string]string{
					"email":    email,
					"password": password,
				})
			},
		).Bind(apis.RequireAuth())
		return se.Next()
	})
}

package main

import (
	"io/fs"
	"path"
	"pb_launcher/ui"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var allowedExts = map[string]bool{
	".html":  true,
	".js":    true,
	".css":   true,
	".png":   true,
	".jpg":   true,
	".jpeg":  true,
	".svg":   true,
	".ico":   true,
	".json":  true,
	".woff":  true,
	".woff2": true,
}

func ServeEmbeddedUI(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.Any("/{path...}", func(e *core.RequestEvent) error {
			reqPath := strings.TrimLeft(path.Clean(e.Request.URL.Path), "/")
			if reqPath == "" || reqPath == "/" {
				reqPath = "index.html"
			}

			info, err := fs.Stat(ui.DistDirFS, reqPath)
			if err != nil || info.IsDir() {
				reqPath = "index.html"
			}
			if ext := path.Ext(reqPath); err == nil {
				if !allowedExts[ext] {
					reqPath = "index.html"
				}
			}
			return e.FileFS(ui.DistDirFS, reqPath)
		})
		return se.Next()
	})
}

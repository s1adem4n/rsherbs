package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"rsherbs/frontend"
	"rsherbs/pkg/labels"

	_ "image/jpeg"
	_ "image/png"
	_ "rsherbs/migrations"
)

func SPAMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code == http.StatusNotFound {
					c.Request().URL.Path = "/"
					return next(c)
				}
			}
		}
		return err
	}
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	var dataDir string
	if isGoRun {
		dataDir = "pb_data"
	} else {
		dataDir = filepath.Join(configDir, "rsherbs")
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
	})

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		router := e.Router
		e.Router.Use(SPAMiddleware)
		e.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, "/_/")
			},
		}))

		subFS := echo.MustSubFS(frontend.Assets, "build")
		router.StaticFS("/", subFS)

		router.GET("/labels", func(c echo.Context) error {
			ids := strings.Split(c.QueryParam("ids"), ",")

			if len(ids) == 0 {
				return c.String(400, "ids required")
			}

			var models []*models.Record

			for _, id := range ids {
				model, err := app.Dao().FindRecordById("plants", id)
				if err != nil {
					return err
				}
				models = append(models, model)
			}

			quantity, err := strconv.Atoi(c.QueryParam("quantity"))
			if err != nil {
				quantity = 1
			}
			width, err := strconv.ParseFloat(c.QueryParam("width"), 64)
			if err != nil {
				width = 80
			}
			print := false
			if c.QueryParam("print") == "true" {
				print = true
			}

			var pdf []byte

			options := labels.GeneratePDFOptions{
				Print:    print,
				Margin:   3,
				Gap:      3,
				Width:    width,
				Height:   width * 9 / 16, // 16 / 9
				Quantity: quantity,
			}

			if len(models) == 1 {
				model := models[0]
				pdf, err = labels.GeneratePDF(labels.LabelMetadata{
					Name:  model.GetString("name"),
					Latin: model.GetString("latin"),
				}, options)
				if err != nil {
					return err
				}
			} else if len(models) > 1 {
				var meta []labels.LabelMetadata
				for _, model := range models {
					meta = append(meta, labels.LabelMetadata{
						Name:  model.GetString("name"),
						Latin: model.GetString("latin"),
					})
				}
				pdf, err = labels.GeneratePDFMultiple(meta, options)
				if err != nil {
					return err
				}
			} else {
				return c.String(400, "no models found")
			}

			return c.Blob(200, "application/pdf", pdf)
		})

		return nil
	})

	shouldAutoServe := !isGoRun && len(os.Args) == 1

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		if shouldAutoServe {
			open("http://localhost:8090/")
			_, err := apis.Serve(app, apis.ServeConfig{
				HttpAddr:        "127.0.0.1:8090",
				ShowStartBanner: false,
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	if shouldAutoServe {
		if err := app.Bootstrap(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}
}

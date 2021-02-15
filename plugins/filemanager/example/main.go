package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	_ "github.com/qtoad/xgo-admin/adapter/gin"
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/sqlite"
	_ "github.com/qtoad/xgo-admin/themes/sword"

	"github.com/gin-gonic/gin"
	"github.com/qtoad/xgo-admin/engine"
	"github.com/qtoad/xgo-admin/modules/config"
	"github.com/qtoad/xgo-admin/modules/language"
	"github.com/qtoad/xgo-admin/plugins/filemanager"
	"github.com/qtoad/xgo-admin/plugins/filemanager/modules/root"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	e := engine.Default()

	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Driver: config.DriverSqlite,
				File:   "./admin.db",
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:           language.EN,
		IndexUrl:           "/fm/def/list",
		Debug:              true,
		AccessAssetsLogOff: true,
		Theme:              "sword",
		Animation: config.PageAnimation{
			Type: "fadeInUp",
		},
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := e.AddConfig(cfg).
		AddPlugins(filemanager.
			NewFileManager(filepath.Join(dir, "book")).
			AddRoot("code", root.Root{Path: filepath.Join(dir, "code"), Title: "Code"}).
			AddRoot("picture", root.Root{Path: filepath.Join(dir, "picture"), Title: "Picture"}),
		).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	go func() {
		_ = r.Run(":9033")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	e.SqliteConnection().Close()
}

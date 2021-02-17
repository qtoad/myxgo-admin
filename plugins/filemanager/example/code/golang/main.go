package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/qtoad/xgo-admin/plugins/filemanager/modules/root"

	_ "github.com/qtoad/xgo-admin/adapter/gin"
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/sqlite"
	_ "github.com/qtoad/xgo-admin/themes/sword"

	"github.com/gin-gonic/gin"
	"github.com/qtoad/xgo-admin/engine"
	"github.com/qtoad/xgo-admin/modules/config"
	"github.com/qtoad/xgo-admin/modules/language"
	"github.com/qtoad/xgo-admin/plugins/filemanager"
)

func main() {
	r := gin.Default()

	//gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard

	e := engine.Default()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := e.AddConfig(config.Config{
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
		Language: language.EN,
		IndexUrl: "/fm/def/list",
		Debug:    true,
		Theme:    "sword",
		Animation: config.PageAnimation{
			Type: "fadeInUp",
		},
	}).
		AddPlugins(filemanager.NewFileManager(dir+"/root1").
			AddRoot("root2", root.Root{Path: dir + "/root2", Title: "root2"}).
			AddRoot("root3", root.Root{Path: dir + "/root3", Title: "root3"}),
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

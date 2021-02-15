package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/qtoad/xgo-admin/adapter/gin"               // web framework adapter
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/sqlite" // sql driver
	_ "github.com/qtoad/xgo-admin/themes/adminlte"           // ui theme

	"github.com/gin-gonic/gin"
	"github.com/qtoad/xgo-admin/demo123/models"
	"github.com/qtoad/xgo-admin/demo123/pages"
	"github.com/qtoad/xgo-admin/demo123/tables"
	"github.com/qtoad/xgo-admin/engine"
	"github.com/qtoad/xgo-admin/template"
	"github.com/qtoad/xgo-admin/template/chartjs"
)

func main() {
	startServer()
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	eng := engine.Default()

	template.AddComp(chartjs.NewChart())

	//cfg := config.Config{
	//	Databases: config.DatabaseList{
	//		"default": {
	//			Host:       "127.0.0.1",
	//			Port:       "3306",
	//			User:       "root",
	//			Pwd:        "root",
	//			Name:       "go-admin",
	//			MaxIdleCon: 50,
	//			MaxOpenCon: 150,
	//			Driver:     db.DriverMysql,
	//		},
	//	},
	//	UrlPrefix: "admin",
	//	IndexUrl:  "/",
	//	Debug:     true,
	//	Language:  language.CN,
	//}

	if err := eng.AddConfigFromJSON("./config.json").
		AddGenerators(tables.Generators).
		AddGenerator("external", tables.GetExternalTable).
		Use(r); err != nil {
		panic(err)
	}

	models.Init(eng.SqliteConnection())

	r.Static("/uploads", "./uploads")

	eng.HTML("GET", "/admin", pages.DashboardPage)
	eng.HTML("GET", "/admin/form", pages.GetFormContent)
	eng.HTML("GET", "/admin/table", pages.GetTableContent)
	eng.HTMLFile("GET", "/admin/hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	srv := &http.Server{
		Addr:    ":9033",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

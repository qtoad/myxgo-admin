package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	_ "github.com/qtoad/myxgo-admin/adapter/gin"
	_ "github.com/qtoad/myxgo-admin/modules/db/drivers/mysql"
	_ "github.com/qtoad/myxgo-admin/themes/sword"

	"github.com/gin-gonic/gin"
	"github.com/qtoad/myxgo-admin/engine"
	"github.com/qtoad/myxgo-admin/examples/datamodel"
	"github.com/qtoad/myxgo-admin/modules/config"
	"github.com/qtoad/myxgo-admin/modules/language"
	"github.com/qtoad/myxgo-admin/plugins/example"
	"github.com/qtoad/myxgo-admin/template"
	"github.com/qtoad/myxgo-admin/template/chartjs"
	"github.com/qtoad/myxgo-admin/themes/adminlte"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.New()

	e := engine.Default()

	cfg := config.Config{
		Env: config.EnvLocal,
		Databases: config.DatabaseList{
			"default": {
				Host:       "127.0.0.1",
				Port:       "8889",
				User:       "root",
				Pwd:        "root",
				Name:       "goadmin",
				MaxIdleCon: 50,
				MaxOpenCon: 150,
				Driver:     config.DriverMysql,

				//Driver: config.DriverSqlite,
				//File:   "../datamodel/admin.db",
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:           language.CN,
		IndexUrl:           "/",
		Debug:              true,
		AccessAssetsLogOff: true,
		Animation: config.PageAnimation{
			Type: "fadeInUp",
		},
		ColorScheme:       adminlte.ColorschemeSkinBlack,
		BootstrapFilePath: "./../datamodel/bootstrap.go",
	}

	template.AddComp(chartjs.NewChart())

	// customize a plugin

	examplePlugin := example.NewExample()

	// load from golang.Plugin
	//
	// examplePlugin := plugins.LoadFromPlugin("../datamodel/example.so")

	// customize the login page
	// example: https://github.com/GoAdminGroup/demo.go-admin.cn/blob/master/main.go#L39
	//
	// template.AddComp("login", datamodel.LoginPage)

	// load config from json file
	//
	// e.AddConfigFromJSON("../datamodel/config.json")

	if err := e.AddConfig(cfg).
		AddGenerators(datamodel.Generators).
		// add generator, first parameter is the url prefix of table when visit.
		// example:
		//
		// "user" => http://localhost:9033/admin/info/user
		//
		AddGenerator("user", datamodel.GetUserTable).
		AddDisplayFilterXssJsFilter().
		AddPlugins(examplePlugin).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	// customize your pages

	e.HTML("GET", "/admin", datamodel.GetContent)

	go func() {
		_ = r.Run(":9033")
		fmt.Println("server port: 9033")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	e.MysqlConnection().Close()
}

package echo

import (
	// add echo adapter
	_ "github.com/qtoad/mygo-admin/adapter/echo"
	"github.com/qtoad/mygo-admin/modules/config"
	"github.com/qtoad/mygo-admin/modules/language"
	"github.com/qtoad/mygo-admin/plugins/admin/modules/table"
	"github.com/qtoad/mygo-admin/themes/adminlte"

	// add mysql driver
	_ "github.com/qtoad/mygo-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/qtoad/mygo-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/qtoad/mygo-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/qtoad/mygo-admin/modules/db/drivers/mssql"
	// add adminlte ui theme
	_ "github.com/GoAdminGroup/themes/adminlte"

	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/qtoad/mygo-admin/engine"
	"github.com/qtoad/mygo-admin/plugins/admin"
	"github.com/qtoad/mygo-admin/plugins/example"
	"github.com/qtoad/mygo-admin/template"
	"github.com/qtoad/mygo-admin/template/chartjs"
	"github.com/qtoad/mygo-admin/tests/tables"
)

func newHandler() http.Handler {
	e := echo.New()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)
	template.AddComp(chartjs.NewChart())

	examplePlugin := example.NewExample()

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(e); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return e
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	e := echo.New()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(gens)

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(config.Config{
		Databases: dbs,
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}).
		AddPlugins(adminPlugin).Use(e); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return e
}

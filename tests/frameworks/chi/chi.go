package chi

import (
	// add chi adapter
	_ "github.com/qtoad/myxgo-admin/adapter/chi"
	"github.com/qtoad/myxgo-admin/modules/config"
	"github.com/qtoad/myxgo-admin/modules/language"
	"github.com/qtoad/myxgo-admin/plugins/admin/modules/table"
	"github.com/qtoad/myxgo-admin/themes/adminlte"

	// add mysql driver
	_ "github.com/qtoad/myxgo-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/qtoad/myxgo-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/qtoad/myxgo-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/qtoad/myxgo-admin/modules/db/drivers/mssql"
	// add adminlte ui theme
	_ "github.com/GoAdminGroup/themes/adminlte"

	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/qtoad/myxgo-admin/engine"
	"github.com/qtoad/myxgo-admin/plugins/admin"
	"github.com/qtoad/myxgo-admin/plugins/example"
	"github.com/qtoad/myxgo-admin/template"
	"github.com/qtoad/myxgo-admin/template/chartjs"
	"github.com/qtoad/myxgo-admin/tests/tables"
)

func newHandler() http.Handler {
	r := chi.NewRouter()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)
	examplePlugin := example.NewExample()
	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(r); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return r
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	r := chi.NewRouter()

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
		AddPlugins(adminPlugin).Use(r); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return r
}

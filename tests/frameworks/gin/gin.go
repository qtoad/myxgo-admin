package gin

import (
	// add gin adapter
	ada "github.com/qtoad/xgo-admin/adapter/gin"
	// add mysql driver
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/qtoad/xgo-admin/modules/db/drivers/mssql"
	// add adminlte ui theme
	_ "github.com/qtoad/xgo-admin/themes/adminlte"

	"io/ioutil"
	"net/http"
	"os"

	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
	"github.com/qtoad/xgo-admin/engine"
	"github.com/qtoad/xgo-admin/modules/config"
	"github.com/qtoad/xgo-admin/modules/language"
	"github.com/qtoad/xgo-admin/plugins/admin/modules/table"
	"github.com/qtoad/xgo-admin/template"
	"github.com/qtoad/xgo-admin/template/chartjs"
	"github.com/qtoad/xgo-admin/tests/tables"
)

func newHandler() http.Handler {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	eng := engine.Default()

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddGenerators(tables.Generators).
		AddGenerator("user", tables.GetUserTable).
		Use(r); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	r.Static("/uploads", "./uploads")

	return r
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	eng := engine.Default()

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
		AddAdapter(new(ada.Gin)).
		AddGenerators(gens).
		Use(r); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	r.Static("/uploads", "./uploads")

	return r
}

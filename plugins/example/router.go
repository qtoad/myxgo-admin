package example

import (
	"github.com/qtoad/myxgo-admin/context"
	"github.com/qtoad/myxgo-admin/modules/auth"
	"github.com/qtoad/myxgo-admin/modules/db"
	"github.com/qtoad/myxgo-admin/modules/service"
)

func (e *Example) initRouter(prefix string, srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(prefix)
	route.GET("/example", auth.Middleware(db.GetConnection(srv)), e.TestHandler)

	return app
}

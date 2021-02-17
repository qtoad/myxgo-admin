package main

import (
	_ "github.com/qtoad/mygo-admin/adapter/gin"
	_ "github.com/qtoad/mygo-admin/modules/db/drivers/sqlite"
	_ "github.com/qtoad/mygo-admin/themes/sword"

	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/qtoad/mygo-admin/context"
	"github.com/qtoad/mygo-admin/engine"
	"github.com/qtoad/mygo-admin/modules/auth"
	"github.com/qtoad/mygo-admin/modules/config"
	"github.com/qtoad/mygo-admin/modules/language"
	"github.com/qtoad/mygo-admin/plugins/admin/models"
	"github.com/qtoad/mygo-admin/plugins/filemanager"
	"github.com/qtoad/mygo-admin/plugins/librarian"
	"github.com/qtoad/mygo-admin/plugins/librarian/modules/theme"
	"github.com/qtoad/mygo-admin/template/types/action"
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
		Language:                      language.EN,
		IndexUrl:                      "/librarian/README",
		Debug:                         true,
		AccessAssetsLogOff:            true,
		HideConfigCenterEntrance:      true,
		HideAppInfoEntrance:           true,
		HideVisitorUserCenterEntrance: true,
		Logo:                          "<b>Li</b>brarian",
		MiniLogo:                      "Li",
		Theme:                         "sword",
		Title:                         "Librarian",
		//ExcludeThemeComponents:        []string{"datatable", "form"},
		//Animation: config.PageAnimation{
		//	Type: "fadeInUp",
		//},
	}

	theme.Set(theme.Config{
		HideNavBar:   true,
		HideMenuIcon: true,
		FixedSidebar: true,
		ChangeTitle:  true,
	})

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	const visitorRoleID = int64(3)

	if err := e.AddConfig(cfg).
		AddNavButtons("Menu", "", action.Jump("/admin/menu")).
		AddNavButtons("Files", "", action.Jump("/admin/fm/def/list")).
		//AddNavButtons("", icon.Pencil, action.Jump("/admin/menu")).
		AddPlugins(librarian.NewLibrarianWithConfig(librarian.Config{
			Path:           filepath.Join(dir, "docs"),
			MenuUserRoleID: visitorRoleID,
			Prefix:         "librarian",
			BuildMenu:      true, // auto build menus
		}), filemanager.NewFileManager(filepath.Join(dir, "docs"))).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	e.Data("GET", "/admin/librarian", func(ctx *context.Context) {
		conn := e.SqliteConnection()
		user := models.User().SetConn(conn).Find(visitorRoleID)
		_ = auth.SetCookie(ctx, user, conn)
		ctx.Redirect("/admin/librarian/README")
	}, true)

	go func() {
		_ = r.Run(":9033")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	e.SqliteConnection().Close()
}

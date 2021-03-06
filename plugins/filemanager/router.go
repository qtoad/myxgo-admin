package filemanager

import (
	"github.com/qtoad/myxgo-admin/context"
	"github.com/qtoad/myxgo-admin/modules/auth"
	"github.com/qtoad/myxgo-admin/modules/service"
)

func (f *FileManager) initRouter(srvList service.List) *context.App {

	app := context.NewApp()
	authRoute := app.Group("/", auth.Middleware(f.Conn))

	authRoute.GET("/", f.guard.Files, f.handler.ListFiles)
	authRoute.GET("/:__prefix/list", f.guard.Files, f.handler.ListFiles)
	authRoute.GET("/:__prefix/download", f.handler.Download)
	authRoute.POST("/:__prefix/upload", f.guard.Upload, f.handler.Upload)
	authRoute.POST("/:__prefix/create/dir/popup", f.handler.CreateDirPopUp)
	authRoute.POST("/:__prefix/create/dir", f.guard.CreateDir, f.handler.CreateDir)
	authRoute.POST("/:__prefix/delete", f.guard.Delete, f.handler.Delete)
	authRoute.POST("/:__prefix/move/popup", f.handler.MovePopup)
	authRoute.POST("/:__prefix/move", f.guard.Move, f.handler.Move)
	authRoute.GET("/:__prefix/preview", f.guard.Preview, f.handler.Preview)
	authRoute.POST("/:__prefix/rename/popup", f.handler.RenamePopUp)
	authRoute.POST("/:__prefix/rename", f.guard.Rename, f.handler.Rename)

	return app
}

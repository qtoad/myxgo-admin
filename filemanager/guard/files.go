package guard

import (
	"github.com/qtoad/xgo-admin/context"
	errors "github.com/qtoad/xgo-admin/filemanager/modules/error"
	"github.com/qtoad/xgo-admin/modules/util"
)

type FilesParam struct {
	*Base
}

func (g *Guardian) Files(ctx *context.Context) {

	relativePath, path, err := g.getPaths(ctx)

	if !util.IsDirectory(path) {
		err = errors.IsNotDir
	}

	ctx.SetUserValue(filesParamKey, &FilesParam{
		Base: &Base{
			Path:     relativePath,
			FullPath: path,
			Error:    err,
			Prefix:   g.GetPrefix(ctx),
		},
	})
	ctx.Next()
}

func GetFilesParam(ctx *context.Context) *FilesParam {
	return ctx.UserValue[filesParamKey].(*FilesParam)
}

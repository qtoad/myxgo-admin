package guard

import (
	"github.com/qtoad/myxgo-admin/context"
	"github.com/qtoad/myxgo-admin/modules/util"
	errors "github.com/qtoad/myxgo-admin/plugins/filemanager/modules/error"
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

package guard

import (
	"github.com/qtoad/xgo-admin/context"
	"github.com/qtoad/xgo-admin/filemanager/modules/constant"
	errors "github.com/qtoad/xgo-admin/filemanager/modules/error"
	"github.com/qtoad/xgo-admin/modules/util"
)

type PreviewParam struct {
	Base
}

func (g *Guardian) Preview(ctx *context.Context) {

	relativePath, path, err := g.getPaths(ctx)

	if !util.IsFile(path) {
		err = errors.IsNotFile
	}

	ctx.SetUserValue(previewParamKey, &PreviewParam{
		Base: Base{
			Path:     relativePath,
			FullPath: path,
			Error:    err,
			Prefix:   ctx.Query(constant.PrefixKey),
		},
	})
	ctx.Next()
}

func GetPreviewParam(ctx *context.Context) *PreviewParam {
	return ctx.UserValue[previewParamKey].(*PreviewParam)
}

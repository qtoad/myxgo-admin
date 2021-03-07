package guard

import (
	"path/filepath"
	"strings"

	"github.com/qtoad/myxgo-admin/context"
	errors "github.com/qtoad/myxgo-admin/plugins/filemanager/modules/error"
	"github.com/qtoad/myxgo-admin/util"
)

type DeleteParam struct {
	Path   string
	Prefix string
	Error  error
	Paths  []string
}

func (g *Guardian) Delete(ctx *context.Context) {

	if !g.permissions.AllowDelete {
		ctx.SetUserValue(deleteParamKey, &DeleteParam{Error: errors.NoPermission})
		ctx.Next()
		return
	}

	var (
		relativePaths   = ctx.FormValue("id")
		relativePathArr = strings.Split(relativePaths, ",")

		paths = make([]string, 0)
	)

	for _, relativePath := range relativePathArr {
		path := filepath.Join(g.roots.GetPathFromPrefix(ctx), relativePath)

		if relativePath == "" || !strings.Contains(path, g.roots.GetPathFromPrefix(ctx)) || !util.FileExist(path) || strings.Contains(path, "..") {
			ctx.SetUserValue(deleteParamKey, &DeleteParam{Error: errors.DirIsNotExist})
			ctx.Next()
			return
		}

		paths = append(paths, path)

	}
	ctx.SetUserValue(deleteParamKey, &DeleteParam{
		Path:   relativePaths,
		Paths:  paths,
		Prefix: g.GetPrefix(ctx),
	})
	ctx.Next()
}

func GetDeleteParam(ctx *context.Context) *DeleteParam {
	return ctx.UserValue[deleteParamKey].(*DeleteParam)
}

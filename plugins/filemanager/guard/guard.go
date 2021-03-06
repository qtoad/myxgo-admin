package guard

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/qtoad/myxgo-admin/context"
	"github.com/qtoad/myxgo-admin/modules/db"
	"github.com/qtoad/myxgo-admin/plugins/filemanager/modules/constant"
	errors "github.com/qtoad/myxgo-admin/plugins/filemanager/modules/error"
	"github.com/qtoad/myxgo-admin/plugins/filemanager/modules/permission"
	"github.com/qtoad/myxgo-admin/plugins/filemanager/modules/root"
	"github.com/qtoad/myxgo-admin/util"
)

type Guardian struct {
	conn        db.Connection
	roots       root.Roots
	permissions permission.Permission
}

func New(roots root.Roots, conn db.Connection, permission permission.Permission) *Guardian {
	return &Guardian{
		roots:       roots,
		conn:        conn,
		permissions: permission,
	}
}

func (g *Guardian) Update(roots root.Roots, permission permission.Permission) {
	g.roots = roots
	g.permissions = permission
}

const (
	filesParamKey     = "files_param"
	uploadParamKey    = "upload_param"
	createDirParamKey = "create_dir_param"
	deleteParamKey    = "delete_param"
	renameParamKey    = "rename_param"
	previewParamKey   = "preview_param"
)

type Base struct {
	Path     string
	Prefix   string
	FullPath string
	Error    error
}

func (g *Guardian) GetPrefix(ctx *context.Context) string {
	prefix := ctx.Query(constant.PrefixKey)
	if prefix == "" {
		if _, ok := g.roots["def"]; ok {
			return "def"
		} else {
			for name := range g.roots {
				return name
			}
		}
	}
	return prefix
}

func (g *Guardian) getPaths(ctx *context.Context) (string, string, error) {
	var (
		err error

		relativePath, _ = url.QueryUnescape(ctx.Query("path"))
		path            = filepath.Join(g.roots.GetPathFromPrefix(ctx), relativePath)
	)
	if !strings.Contains(path, g.roots.GetPathFromPrefix(ctx)) {
		err = errors.DirIsNotExist
	}

	if !util.FileExist(path) {
		err = errors.DirIsNotExist
	}

	return relativePath, path, err
}

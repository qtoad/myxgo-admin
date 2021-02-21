package guard

import (
	"path/filepath"
	"strings"

	"github.com/qtoad/mygo-admin/context"
	"github.com/qtoad/mygo-admin/modules/config"
	"github.com/qtoad/mygo-admin/modules/db"
	"github.com/qtoad/mygo-admin/plugins/librarian/modules/constant"
	"github.com/qtoad/mygo-admin/plugins/librarian/modules/root"
)

type Guardian struct {
	conn   db.Connection
	roots  *root.Roots
	prefix string
}

func New(roots *root.Roots, conn db.Connection, prefix string) *Guardian {
	return &Guardian{
		roots:  roots,
		conn:   conn,
		prefix: prefix,
	}
}

const (
	viewParamKey = "view_param"
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
		return "def"
	}
	return prefix
}

func (g *Guardian) Update(root *root.Roots) {
	g.roots = root
}

func (g *Guardian) getPaths(ctx *context.Context) (string, string, error) {
	var (
		err          error
		relativePath = strings.Replace(ctx.Path(), config.Url("/"+g.prefix), "", -1) + ".md"
		path         = filepath.Join(g.roots.GetPathFromPrefix(ctx), relativePath)
	)

	return relativePath, path, err
}

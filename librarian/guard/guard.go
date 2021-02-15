package guard

import (
	"path/filepath"
	"strings"

	"github.com/qtoad/xgo-admin/context"
	"github.com/qtoad/xgo-admin/librarian/modules/constant"
	"github.com/qtoad/xgo-admin/librarian/modules/root"
	"github.com/qtoad/xgo-admin/modules/config"
	"github.com/qtoad/xgo-admin/modules/db"
)

type Guardian struct {
	conn   db.Connection
	roots  *root.Roots
	prefix string
}

func New(r *root.Roots, c db.Connection, p string) *Guardian {
	return &Guardian{
		roots:  r,
		conn:   c,
		prefix: p,
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

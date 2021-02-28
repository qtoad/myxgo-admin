package controller

import (
	"github.com/qtoad/myxgo-admin/context"
	"github.com/qtoad/myxgo-admin/plugins/librarian/modules/constant"
	"github.com/qtoad/myxgo-admin/plugins/librarian/modules/root"
	"github.com/qtoad/myxgo-admin/template"
	"github.com/qtoad/myxgo-admin/template/types"
)

type Handler struct {
	roots *root.Roots
	theme string

	HTML func(ctx *context.Context, panel types.Panel, options ...template.ExecuteOptions)
}

func NewHandler(root *root.Roots, theme string) *Handler {
	return &Handler{
		roots: root,
		theme: theme,
	}
}

func (h *Handler) Prefix(ctx *context.Context) string {
	prefix := ctx.Query(constant.PrefixKey)
	if prefix == "" {
		return "def"
	}
	return prefix
}

func (h *Handler) Update(root *root.Roots, theme string) {
	h.roots = root
	h.theme = theme
}

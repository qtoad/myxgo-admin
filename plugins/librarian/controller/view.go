package controller

import (
	"io/ioutil"

	"github.com/qtoad/xgo-admin/template"

	"github.com/qtoad/xgo-admin/context"
	"github.com/qtoad/xgo-admin/modules/util"
	"github.com/qtoad/xgo-admin/plugins/librarian/guard"
	"github.com/qtoad/xgo-admin/plugins/librarian/modules/theme"
	"github.com/qtoad/xgo-admin/template/types"
)

func (h *Handler) View(ctx *context.Context) {

	param := guard.GetViewParam(ctx)

	content, err := ioutil.ReadFile(param.FullPath)

	if err != nil {
		panic(err)
	}

	var md = /*blackfriday.Run(*/ util.NormalizeEOL(content) //)
	h.HTML(ctx, types.Panel{
		Content: theme.Get(h.theme).HTML(md),
		CSS:     theme.Get(h.theme).CSS(),
		JS:      theme.Get(h.theme).JS(),
	}, template.ExecuteOptions{
		NoCompress: true,
		Animation:  true,
	})
}

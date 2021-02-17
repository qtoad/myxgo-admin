package controller

import (
	"io/ioutil"

	"github.com/qtoad/mygo-admin/template"

	"github.com/qtoad/mygo-admin/context"
	"github.com/qtoad/mygo-admin/modules/util"
	"github.com/qtoad/mygo-admin/plugins/librarian/guard"
	"github.com/qtoad/mygo-admin/plugins/librarian/modules/theme"
	"github.com/qtoad/mygo-admin/template/types"
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

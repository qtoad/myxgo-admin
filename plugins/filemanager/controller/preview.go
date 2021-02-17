package controller

import (
	"github.com/qtoad/mygo-admin/context"
	"github.com/qtoad/mygo-admin/plugins/filemanager/guard"
	"github.com/qtoad/mygo-admin/plugins/filemanager/previewer"
)

func (h *Handler) Preview(ctx *context.Context) {
	param := guard.GetPreviewParam(ctx)
	if param.Error != nil {
		h.preview(ctx, "", param.Path, param.FullPath, param.Error)
		return
	}
	content, err := previewer.Preview(param.FullPath)
	h.preview(ctx, content, param.Path, param.FullPath, err)
}

package controller

import (
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"

	"github.com/qtoad/xgo-admin/context"
	"github.com/qtoad/xgo-admin/filemanager/guard"
	"github.com/qtoad/xgo-admin/filemanager/models"
	"github.com/qtoad/xgo-admin/modules/util"
)

func (h *Handler) ListFiles(ctx *context.Context) {

	var (
		param      = guard.GetFilesParam(ctx)
		filesOfDir = make(models.Files, 0)
	)

	if param.Error != nil {
		h.table(ctx, filesOfDir, param.Error)
		return
	}

	fileInfos, err := ioutil.ReadDir(filepath.FromSlash(param.FullPath))

	if err != nil {
		h.table(ctx, filesOfDir, err)
		return
	}

	for _, fileInfo := range fileInfos {

		if util.IsHiddenFile(fileInfo.Name()) {
			continue
		}

		file := models.File{
			IsDirectory:  fileInfo.IsDir(),
			Name:         fileInfo.Name(),
			Size:         int(fileInfo.Size()),
			Extension:    strings.TrimLeft(filepath.Ext(fileInfo.Name()), "."),
			Path:         filepath.Join(param.Path, fileInfo.Name()),
			Mime:         mime.TypeByExtension(filepath.Ext(fileInfo.Name())),
			LastModified: fileInfo.ModTime().Unix(),
		}

		filesOfDir = append(filesOfDir, file)
	}

	h.table(ctx, filesOfDir, nil)
	return
}

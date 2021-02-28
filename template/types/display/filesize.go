package display

import (
	"strconv"

	"github.com/qtoad/myxgo-admin/modules/util"
	"github.com/qtoad/myxgo-admin/template/types"
)

type FileSize struct {
	types.BaseDisplayFnGenerator
}

func init() {
	types.RegisterDisplayFnGenerator("filesize", new(FileSize))
}

func (f *FileSize) Get(args ...interface{}) types.FieldFilterFn {
	return func(value types.FieldModel) interface{} {
		size, _ := strconv.ParseUint(value.Value, 10, 64)
		return util.FileSize(size)
	}
}

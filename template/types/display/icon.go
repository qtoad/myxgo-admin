package display

import (
	"github.com/qtoad/myxgo-admin/template/icon"
	"github.com/qtoad/myxgo-admin/template/types"
)

type Icon struct {
	types.BaseDisplayFuncGenerator
}

func init() {
	types.RegisterDisplayFnGenerator("icon", new(Icon))
}

func (i *Icon) Get(args ...interface{}) types.FieldFilterFunc {
	return func(value types.FieldModel) interface{} {
		icons := args[0].(map[string]string)
		defaultIcon := ""
		if len(args) > 1 {
			defaultIcon = args[1].(string)
		}
		for k, iconClass := range icons {
			if k == value.Value {
				return icon.Icon(iconClass)
			}
		}
		if defaultIcon != "" {
			return icon.Icon(defaultIcon)
		}
		return value.Value
	}
}

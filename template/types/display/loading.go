package display

import (
	"html/template"

	"github.com/qtoad/myxgo-admin/template/types"
)

type Loading struct {
	types.BaseDisplayFuncGenerator
}

func init() {
	types.RegisterDisplayFnGenerator("loading", new(Loading))
}

func (l *Loading) Get(args ...interface{}) types.FieldFilterFunc {
	return func(value types.FieldModel) interface{} {
		param := args[0].([]string)

		for i := 0; i < len(param); i++ {
			if value.Value == param[i] {
				return template.HTML(`<i class="fa fa-refresh fa-spin text-primary"></i>`)
			}
		}

		return value.Value
	}
}

package display

import (
	"html/template"

	"github.com/qtoad/myxgo-admin/template/types"
)

type Copyable struct {
	types.BaseDisplayFuncGenerator
}

func init() {
	types.RegisterDisplayFnGenerator("copyable", new(Copyable))
}

func (c *Copyable) Get(args ...interface{}) types.FieldFilterFunc {
	return func(value types.FieldModel) interface{} {
		return template.HTML(`
<a href="javascript:void(0);" class="grid-column-copyable text-muted" data-content="` + value.Value + `" 
title="Copied!" data-placement="bottom">
<i class="fa fa-copy"></i>
</a>&nbsp;` + value.Value + `
`)
	}
}

func (c *Copyable) JS() template.HTML {
	return template.HTML(`
$('body').on('click','.grid-column-copyable',(function (e) {
	var content = $(this).data('content');
	
	var temp = $('<input>');
	
	$("body").append(temp);
	temp.val(content).select();
	document.execCommand("copy");
	temp.remove();
	
	$(this).tooltip('show');
}));
`)
}

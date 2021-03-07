package components

import (
	"github.com/qtoad/myxgo-admin/util"
	"html/template"

	"github.com/qtoad/myxgo-admin/template/types"
)

type ImgAttribute struct {
	Name     string
	Width    string
	Height   string
	Uuid     string
	HasModal bool
	Src      template.URL
	types.Attribute
}

func (compo *ImgAttribute) SetWidth(value string) types.ImgAttribute {
	compo.Width = value
	return compo
}

func (compo *ImgAttribute) SetHeight(value string) types.ImgAttribute {
	compo.Height = value
	return compo
}

func (compo *ImgAttribute) WithModal() types.ImgAttribute {
	compo.HasModal = true
	compo.Uuid = util.NewUuid()
	return compo
}

func (compo *ImgAttribute) SetSrc(value template.HTML) types.ImgAttribute {
	compo.Src = template.URL(value)
	return compo
}

func (compo *ImgAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, compo.Separation, *compo, "image")
}

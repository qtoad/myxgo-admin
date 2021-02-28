package sword

import (
	"strings"

	"github.com/gobuffalo/packr/v2"
	adminTemplate "github.com/qtoad/myxgo-admin/template"
	"github.com/qtoad/myxgo-admin/template/components"
	"github.com/qtoad/myxgo-admin/template/types"
	"github.com/qtoad/myxgo-admin/themes/common"
	"github.com/qtoad/myxgo-admin/themes/sword/resource"
)

type Theme struct {
	ThemeName string
	components.Base
	*common.BaseTheme
}

var Sword = Theme{
	ThemeName: "sword",
	Base: components.Base{
		Attribute: types.Attribute{
			TemplateList: TemplateList,
		},
	},
	BaseTheme: &common.BaseTheme{
		AssetPaths:   resource.AssetPaths,
		TemplateList: TemplateList,
	},
}

func init() {
	adminTemplate.Add("sword", &Sword)
}

func Get() *Theme {
	return &Sword
}

func (t *Theme) Name() string {
	return t.ThemeName
}

func (t *Theme) GetTmplList() map[string]string {
	return TemplateList
}

func (t *Theme) GetAsset(path string) ([]byte, error) {
	path = strings.Replace(path, "/assets/dist", "", -1)
	box := packr.New("sword", "./resource/assets/dist")
	return box.Find(path)
}

func (t *Theme) GetAssetList() []string {
	return resource.AssetsList
}

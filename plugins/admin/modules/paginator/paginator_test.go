package paginator

import (
	"testing"

	"github.com/qtoad/xgo-admin/modules/config"
	"github.com/qtoad/xgo-admin/plugins/admin/modules/parameter"
	_ "github.com/qtoad/xgo-admin/themes/sword"
)

func TestGet(t *testing.T) {
	config.Initialize(&config.Config{Theme: "sword"})
	Get(Config{
		Size:         105,
		Param:        parameter.BaseParam().SetPage("7"),
		PageSizeList: []string{"10", "20", "50", "100"},
	})
}

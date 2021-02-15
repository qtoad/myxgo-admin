package gin

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/qtoad/xgo-admin/tests/common"
)

func TestGin(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(newHandler()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}))
}

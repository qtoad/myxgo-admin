package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/qtoad/myxgo-admin/modules/config"
	"github.com/qtoad/myxgo-admin/tests"
	"github.com/qtoad/myxgo-admin/tests/common"
	"github.com/qtoad/myxgo-admin/tests/frameworks/gin"
	"github.com/qtoad/myxgo-admin/tests/tables"
	"github.com/qtoad/myxgo-admin/tests/web"
)

// Black box testing
func TestDemoBlackBox(t *testing.T) {
	tests.BlackBoxTestSuit(t, gin.NewHandler, config.DatabaseList{
		"default": {
			Driver: config.DriverSqlite,
			File:   "./admin.db",
		},
	}, tables.Generators, func(cfg config.DatabaseList) {
		// Data cleaner of the framework
		tests.Cleaner(cfg)
		// Clean your own data:
		// ...
	}, func(e *httpexpect.Expect) {
		// Test cases of the framework
		common.Test(e)
		// Write your own API test, for example:
		// More usages: https://github.com/gavv/httpexpect
		e.GET("/admin").Expect().Status(http.StatusOK)
	})
}

// User acceptance testing
func TestDemoUserAcceptance(t *testing.T) {
	web.UserAcceptanceTestSuit(t, func(t *testing.T, page *web.Page) {
		// Write test case base on chromedriver, for example:
		// More usages: https://github.com/sclevine/agouti
		page.NavigateTo("http://127.0.0.1:9033/admin")
		page.Contain("username")
		page.Click("")
	}, func(quit chan struct{}) {
		// start the server:
		// ....
	}, true) // if local parameter is true, it will not be headless, and window not close when finishing tests.
}

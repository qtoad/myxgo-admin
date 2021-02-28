package main

import (
	"github.com/qtoad/myxgo-admin/version"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetLatestVersion(t *testing.T) {
	assert.Equal(t, getLatestVersion(), version.Version())
}

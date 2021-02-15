package language

import (
	"html/template"

	"github.com/qtoad/xgo-admin/modules/language"
)

func Get(key string) string {
	return language.GetWithScope(key, "librarian")
}

func GetHTML(key string) template.HTML {
	return template.HTML(language.GetWithScope(key, "librarian"))
}

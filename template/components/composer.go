package components

import (
	"bytes"
	"html/template"

	"github.com/qtoad/myxgo-admin/modules/config"
	"github.com/qtoad/myxgo-admin/modules/logger"
	"github.com/qtoad/myxgo-admin/modules/util"
	template2 "github.com/qtoad/myxgo-admin/template"
)

func ComposeHtml(temList map[string]string, separation bool, compo interface{}, templateName ...string) template.HTML {

	tmplName := ""
	if len(templateName) > 0 {
		tmplName = templateName[0] + " "
	}

	var (
		tmpl *template.Template
		err  error
	)

	if separation {
		files := make([]string, 0)
		root := config.GetAssetRootPath() + "pages/"
		for _, v := range templateName {
			files = append(files, root+temList["components/"+v]+".tmpl")
		}
		tmpl, err = template.New("comp").Funcs(template2.DefaultFuncMap).ParseFiles(files...)
	} else {
		var text = ""
		for _, v := range templateName {
			text += temList["components/"+v]
		}
		tmpl, err = template.New("comp").Funcs(template2.DefaultFuncMap).Parse(text)
	}

	if err != nil {
		logger.Panic(tmplName + "ComposeHtml Error:" + err.Error())
		return ""
	}
	buf := new(bytes.Buffer)

	defineName := util.ReplaceAll(templateName[0], "table/", "", "form/", "")

	err = tmpl.ExecuteTemplate(buf, defineName, compo)
	if err != nil {
		logger.Error(tmplName+" ComposeHtml Error:", err)
	}
	return template.HTML(buf.String())
}

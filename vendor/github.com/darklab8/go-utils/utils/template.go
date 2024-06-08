package utils

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/darklab8/go-utils/utils/utils_logus"
	"github.com/darklab8/go-utils/utils/utils_types"
)

func TmpRender(templateRef *template.Template, data interface{}) string {
	header := bytes.Buffer{}
	err := templateRef.Execute(&header, data)
	utils_logus.Log.CheckPanic(err, "failed to render template")
	return header.String()
}

func TmpInit(content utils_types.TemplateExpression) *template.Template {
	funcs := map[string]any{
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix}

	var err error
	templateRef, err := template.New("test").Funcs(funcs).Parse(string(content))
	utils_logus.Log.CheckPanic(err, "failed to init template")
	return templateRef
}

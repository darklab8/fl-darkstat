package web

import (
	"fmt"
	"io"
	"text/template"

	"github.com/darklab8/fl-darkstat/darkcore/settings/logus"
)

var (
	redirect_template *template.Template
)

func init() {
	var err error
	redirect_template, err = template.New("foo").Parse(`
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8"/>
		<title>{ "darkstat oauth" }</title>
		<meta http-equiv="refresh" content="{{.SiteUrl}}"/>
	</head>
	<body>
		{{.Msg}}
	</body>
</html>
`)
	logus.Log.CheckPanic(err, "failed to parse template")
}

type RedirectArgs struct {
	Msg     string
	SiteUrl string
}

func RedirectPageRender(msg string, site_url string, wr io.Writer) error {
	return redirect_template.Execute(wr, RedirectArgs{Msg: msg, SiteUrl: fmt.Sprintf("3; url=%s", site_url)})
}

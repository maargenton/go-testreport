package template

import (
	"bytes"
	"strings"
	"text/template"
)

type Template = template.Template

func New(name string) *Template {
	tmpl := template.New("report")
	tmpl.Funcs(map[string]interface{}{
		"indent":    indent,
		"render":    renderFunc(tmpl),
		"codeblock": codeblock,
	})

	return tmpl
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func renderFunc(tmpl *template.Template) func(name string, data interface{}) (string, error) {
	return func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func codeblock() string {
	return "```"
}

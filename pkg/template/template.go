package template

import (
	"bytes"
	"strings"
	"text/template"
)

// Template is a alias of the same type from `test/template` package
type Template = template.Template

// New allocates a new template with the give name, and attaches an additional
// set of built-in function to help with partial rendering and indentation.
func New(name string) *Template {
	tmpl := template.New(name)
	tmpl.Funcs(map[string]interface{}{
		"indent":    indent,
		"header":    header(0),
		"render":    renderFunc(tmpl),
		"codeblock": codeblock,
	})

	return tmpl
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func header(offset int) func(level int) string {
	return func(level int) string {
		return strings.Repeat("#", level+offset)
	}
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

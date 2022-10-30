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
func New(name string, values map[string]interface{}) *Template {
	var tmpl = template.New(name)

	var headerShift = 0
	if v, ok := values["HeaderShift"]; ok {
		if vv, ok := v.(int); ok {
			headerShift = vv
		}

	}
	var funcs = map[string]interface{}{
		"indent":    indent,
		"header":    header(headerShift),
		"render":    renderFunc(tmpl),
		"codeblock": codeblock,
	}
	for k, v := range values {
		funcs[k] = wrapValue(v)
	}
	tmpl.Funcs(funcs)

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

func wrapValue(v interface{}) func() interface{} {
	return func() interface{} {
		return v
	}
}

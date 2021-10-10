package json_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/maargenton/go-fileutils"
	"github.com/maargenton/go-testpredicate/pkg/bdd"
	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/verify"
	"gopkg.in/yaml.v3"

	"github.com/maargenton/go-testreport/pkg/json"
	"github.com/maargenton/go-testreport/pkg/test"
)

func TestLoad(t *testing.T) {
	bdd.Given(t, "a file containing json formated test output", func(t *bdd.T) {
		filename := "./testdata/test-output.json"

		t.When("doing something", func(t *bdd.T) {
			var pkgs []test.Package
			err := fileutils.ReadFile(filename, func(r io.Reader) (err error) {
				pkgs, err = json.Load(r)
				return err
			})

			fileutils.WriteFile("./testdata/test-output.yaml", func(w io.Writer) error {
				e := yaml.NewEncoder(w)
				return e.Encode(pkgs)
			})

			tmpl := template.New("report")
			tmpl.Funcs(map[string]interface{}{
				"indent":    indent,
				"render":    renderFunc(tmpl),
				"codeblock": codeblock,
			})
			tmpl.Parse(tmplSrc)

			require.That(t, err).IsError(nil)
			require.That(t, tmpl).IsNotNil()

			err = fileutils.WriteFile("./testdata/test-output.md", func(w io.Writer) error {
				return tmpl.Execute(w, pkgs)
			})
			verify.That(t, err).IsError(nil)

			t.Then("something happens", func(t *bdd.T) {
				verify.That(t, err).IsError(nil)
				verify.That(t, pkgs).Length().Eq(6)
			})
		})
	})
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

var tmplSrc = `
# Test report

{{  range . -}}
## {{ .Name }}
{{    range .Tests }}
- {{ .Name -}}
{{      range .LeafTests }}
  - {{ .PartialName 1 -}}
{{      end }}
{{    end }}
{{- end }}

---

{{ range . -}}
{{- template "package" . -}}
{{- end }}


{{- define "package" -}}
{{-   if gt (len .Tests) 0 -}}
## {{ .Name }}

Coverage: {{ .Coverage }}%

{{      range .Tests -}}
{{-       render "test" . -}}
{{-     end }}
{{    end -}}
{{- end -}}


{{- define "test" }}
- {{ render "outcome" . }} {{ .Name }}
{{-   render "output" . | indent 2 -}}
{{-   range .SubTests -}}
{{-     render "test" . | indent 2 -}}
{{-   end -}}
{{- end -}}


{{- define "outcome" }}
{{-   if .Failure -}}❌{{- else -}}✅{{- end -}}
{{- end -}}


{{- define "output" }}
{{-   if gt (len .Output) 0 }}
{{      codeblock -}}
{{      range .Output }}
{{        .  }}
{{-     end }}
{{      codeblock }}
{{-   end -}}
{{- end -}}

`

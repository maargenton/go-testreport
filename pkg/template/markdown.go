package template

// {{  range . -}}
// ## {{ .Name }}
// {{    range .Tests }}
// - {{ .Name -}}
// {{      range .LeafTests }}
//   - {{ .PartialName 1 -}}
// {{      end }}
// {{    end }}
// {{- end }}

// ---

// MarkdownTemplate contains the definition for the built-in template used to
// generate Markdown report
var MarkdownTemplate = `
# Test report
{{ range . -}}
{{- template "package" . -}}
{{- end }}


{{- define "package" -}}
{{-   if gt (len .Tests) 0 }}

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
{{-   range .Tests -}}
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

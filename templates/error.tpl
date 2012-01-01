{{ define "error.contents" }}
  <div id="error-page" class="box 500">
    <h2>{{ .StatusCode }} - {{ httpStatusText .StatusCode }}</h2>
    <p>The request triggered an error.</p>
    <pre>{{ .Error.Error }}</pre>
    <pre>{{ .Stack }}</pre>
  </div>
{{ end }}

{{ define "error" }}
{{ template "base.header" .Common }}
{{ template "error.contents" .Contents }}
{{ template "base.footer" .Common }}
{{ end }}

{{/* vim: set ts=2 sw=2 et tw=0 ft=gotplhtml: */}}

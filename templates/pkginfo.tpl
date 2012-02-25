{{ define "pkginfo.metadata" }}
    <tr>
      <th>Architecture:</th>
      <td>{{ .Architecture }}</td>
    </tr>
    <tr>
      <th>Repository:</th>
      <td><a href="/pkglist?repo={{ .DB.Name }}">{{ .DB.Name }}</a></td>
    </tr>
    <tr>
      <th>Groups:</th>
      <td>{{ range .Groups.Slice }}{{ . }} {{ end }}</td>
    </tr>
    <tr>
      <th>Licenses:</th>
      <td>{{ range .Licenses.Slice }}{{ . }} {{ end }}</td>
    </tr>
    <tr>
      <th>Description:</th>
      <td>{{ .Description }}</td>
    </tr>
    <tr>
      <th>Upstream URL:</th>
      <td><a href="{{ .URL }}">{{ .URL }}</a></td>
    </tr>
    <tr>
      <th>Packager:</th>
      <td>{{ .Packager }}</td>
    </tr>
    <tr>
      <th>Installed size:</th>
      <td>{{ .ISize }} bytes</td>
    </tr>
    <tr>
      <th>Build date:</th>
      <td>{{ .BuildDate }}</td>
    </tr>
    {{ if isLocal . }}
    <tr>
      <th>Install date:</th>
      <td>{{ .InstallDate }}</td>
    </tr>
    <tr>
      <th>Install reason:</th>
      <td>{{ .Reason }}</td>
    </tr>
    {{ end }}
{{ end }}

{{ define "pkginfo.tables" }}
{{ $repo := .Repo }}
<div id="pkgdeps" class="listing">
  <h3 title="Dependencies">Dependencies</h3>
  <ul>
  {{ range $dep := .Package.Depends.Slice }}
  <li><a href="/info?pkg={{ $dep.Name }}&db={{ $repo }}">{{ $dep }}</a></li>
  {{ end }}
  </ul>
</div>
<div id="pkgreqs" class="listing">
  <h3 title="Required by">Required by</h3>
  <ul>
  {{ range $pkgname = .Package.ComputeRequiredBy }}
  <li><a href="/info?pkg={{ $pkgname }}&db={{ $repo }}">{{ $pkgname }}</a></li>
  {{ end }}
  </ul>
  </ul>
</div>
<div id="pkgfiles" class="listing">
  <h3 title="Files">Files listing</h3>
  {{ if .Package.Files }}
  <ul>
  {{ range $fileinfo := .Package.Files }}
  <li><a href="/file/?path=/{{ $fileinfo.Name }}">/{{ $fileinfo.Name }}</a></li>
  {{ end }}
  </ul>
  {{ else }}
  <p>No information available.</p>
  {{ end }}
</div>
{{ end }}

{{ define "pkginfo.contents" }}
<div id="pkgdetails" class="box">
  <h2>Package information: {{ .Package.Name }} {{ .Package.Version }}</h2>
  {{/*
  <div id="detailslinks" class="listing">
    {% if actions %}
    <div id="actionlist">
      <h4>Package actions</h4>
      <ul class="small">
        {% for action, desc in actions %}
        <li><a href="/info?pkg={{ pkg.name }}&db={{ repo }}&action={{ action }}">{{ desc }}</a></li>
        {% endfor %}
      </ul>
    </div>
    {% endif %}
  </div>
  */}}
  <table id="pkginfo">
   {{ template "pkginfo.metadata" $.Package }}
  </table>

  <div id="metadata">
   {{ template "pkginfo.tables" $ }}
  </div>
</div>
{{ end }}

{{ define "pkginfo" }}
{{ template "base.header" .Common }}
{{ template "pkginfo.contents" .Contents }}
{{ template "base.footer" .Common }}
{{ end }}

{{/* vim: set ts=2 sw=2 et tw=0 ft=gotplhtml: */}}

{{ define "repolist.contents" }}
<div class="box">
  <h2>Available repositories</h2>
  <table class="results">
    <thead>
      <tr>
        <th>Name</th>
        <th>URL</th>
        <th>Packages</th>
      </tr>
    </thead>
    <tbody>
      {{ range $db := $.Repos }}
      <tr>
        <td><a href="/pkglist?repo={{ $db.Name }}">{{ $db.Name }}</a></td>
        <td>{{ range $db.Servers }}<a href="{{ $ }}">{{ $ }}</a><br/>{{ end }}</td>
        <td>{{ len $db.PkgCache.Slice }} packages</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</div>
{{ end }}

{{ define "repolist" }}
{{ template "base.header" .Common }}
{{ template "repolist.contents" .Contents }}
{{ template "base.footer" .Common }}
{{ end }}

{{/* vim: set ts=2 sw=2 et tw=0 ft=gotplhtml: */}}

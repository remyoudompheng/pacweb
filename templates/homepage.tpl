{{ define "homepage.contents" }}
<div id="content-left-wrapper">
  <div id="content-left">

<div id="pkg-updates" class="widget box">
  <h3>Recent packages</h3>
  <table>
    <thead>
      <tr>
        <th>Package</th><th>Repository</th><th>Install status</th>
      </tr>
    </thead>
    {{ range $pkg := $.Latest }}
    <tr>
      <td class="pkg-name"><a href="/info?pkg={{ $pkg.Name }}&db={{ $pkg.DB.Name }}">{{ $pkg.Name }} {{ $pkg.Version }}</a></td>
      <td><a href="/pkglist?repo={{ $pkg.DB.Name }}">{{ $pkg.DB.Name }}</a></td>
      <td>{{ installStatus $pkg }}</td>
    </tr>
    {{ end }}
  </table>
</div>

<div id="pkg-outofdate" class="widget box">
  <h3>Packages that can be upgraded</h3>
  <table>
    <thead>
      <tr><th>Package name</th><th>Installed version</th><th>Available version</th></tr>
    </thead>
    {{ range $name, $vers := $.Outdated }}
    <tr>
      <td class="pkg-name">{{ $name }}</td>
      <td><a href="/info?db=local&pkg={{ $name }}">{{ index $vers 0 }}</a></td>
      <td><a href="/info?db={{ index $vers 2 }}&pkg={{ $name }}">{{ index $vers 1 }} ({{ index $vers 2 }})</a></td>
    </tr>
    {{ end }}
  </table>
</div>

  </div>
</div>
<div id="content-right">
  <div id="nav-sidebar" class="widget">
    <ul>
      <li><a href="/?action=update">Update databases</a></li>
    </ul>
  </div>
</div>
{{ end }}

{{ define "homepage" }}
{{ template "base.header" .Common }}
{{ template "homepage.contents" .Contents }}
{{ template "base.footer" .Common }}
{{ end }}

{{/* vim: set ts=2 sw=2 et tw=0 ft=gotplhtml: */}}

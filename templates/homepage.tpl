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
    Not implemented.
    {{/*
    <tr>
      <td class="pkg-name"><a href="/info?pkg={{ name }}&db={{ repo }}">{{ name }} {{ version }}</a></td>
      <td><a href="/pkglist?repo={{ repo }}">{{ repo }}</a></td>
      <td>{{ state }}</td>
    </tr>
    */}}
  </table>
</div>

<div id="pkg-outofdate" class="widget box">
  <h3>Packages that can be upgraded</h3>
  <table>
    <thead>
      <tr><th>Package name</th><th>Installed version</th><th>Available version</th></tr>
    </thead>
    Not implemented.
    {{/*
    <tr>
      <td class="pkg-name">{{ name }}</td>
      <td><a href="/info?db=local&pkg={{ name }}">{{ local_ver }}</a></td>
      <td><a href="/info?db={{ repo }}&pkg={{ name }}">{{ "%s (%s)" % (repo_ver, repo) }}</a></td>
    </tr>
    */}}
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

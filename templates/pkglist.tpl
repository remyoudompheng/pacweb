{{ define "pkglist.contents" }}
<script type="text/javascript">
$(document).ready(function() {
    $(".results").tablesorter({widgets: ['zebra'], sortList: [[0,0]]});
});
</script>
<div class="box">
  <h2>{{ $.Title }}</h2>
  <table class="results">
    <thead>
      <tr>
        <th>Name</th>
        <th>Version</th>
        <th>Description</th>
        <th>Install Size</th>
      </tr>
    </thead>
    <tbody>
      {{ $repo := .Repo }}
      {{ range $idx, $pkg := .Packages }}
      <tr class="{{ parity $idx }}">
        <td><a href="/info?pkg={{ $pkg.Name }}&db={{ $repo }}">{{ $pkg.Name }}</a></td>
        <td>{{ $pkg.Version }}</td>
        <td>{{ $pkg.Description }}</td>
        <td>{{ humanSize $pkg.ISize }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</div>
{{ end }}

{{ define "pkglist" }}
{{ template "base.header" .Common }}
{{ template "pkglist.contents" .Contents }}
{{ template "base.footer" .Common }}
{{ end }}

{{/* vim: set ts=2 sw=2 et tw=0 ft=gotplhtml: */}}

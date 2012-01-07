{{ define "base.header" }}
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
  <head>
    <title>PacWeb: a package management front-end</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <link rel="stylesheet" type="text/css" href="/static/archweb.css" media="screen, projection" />
    <script type="text/javascript" src="/static/jquery-1.4.4.min.js"></script>
    <script type="text/javascript" src="/static/jquery.tablesorter.min.js"></script>
    <script type="text/javascript" src="/static/pacweb.js"></script>
  </head>
  <body>
    <div id="archnavbar" class="{% block navbarclass %}anb-home{% endblock %}">
        <div id="archnavbarlogo"><h1><a href="/" title="Return to the main page">Arch Linux</a></h1></div>
        <div id="archnavbarmenu">
            <ul id="archnavbarlist">
              <li id="anb-home"><a href="/" title="Home page">Home</a></li>
              <li id="anb-local"><a href="/pkglist" title="Local packages">Installed packages</a></li>
              <li id="anb-sync"><a href="/repolist" title="Repositories">Repositories</a></li>
            </ul>
        </div>
    </div><!-- #archnavbar -->
    <div id="content">
      <h1>Pacweb: a web interface to Arch Linux packages</h1>
{{ end }}

{{ define "base.footer" }}
      <div id="footer">
        <p>Copyright © 2011 <a href="mailto:remy@archlinux.org">Rémy Oudompheng</a>.</p>
        <p>Style and decorations from the ArchWeb project
        (Copyright © 2002-2011 <a href="mailto:jvinet@zeroflux.org"
          title="Contact Judd Vinet">Judd Vinet</a> and <a href="mailto:aaron@archlinux.org"
          title="Contact Aaron Griffin">Aaron Griffin</a>).</p>

        <p>The Arch Linux name and logo are recognized
        <a href="https://wiki.archlinux.org/index.php/DeveloperWiki:TrademarkPolicy"
          title="Arch Linux Trademark Policy">trademarks</a>. Some rights reserved.</p>

        <p>The registered trademark Linux® is used pursuant to a sublicense from LMI,
        the exclusive licensee of Linus Torvalds, owner of the mark on a world-wide basis.</p>
      </div><!-- #footer -->
    </div>
  </body>
</html>
{{ end }}

{{ define "base.sysmessage" }}
  {{ if .SysMessage }}
    <div id="sys-message">
      <p>{{ .SysMessage }}</p>
    </div>
  {{ end }}
{{ end }}

{{ define "base.transactionInfo" }}
<div id="transdetails" class="box">
  {{ if .Add }}
  <div id="pkgadd" class="listing">
    <h3 title="Added packages">Added packages</h3>
    <ul>
      {{ range .Add }}
      <li><a href="/info?pkg={{ .Name }}&db={{ .Db.Name }}">{{ .Name }}-{{ .Version }}</a></li>
      {{ end }}
    </ul>
  </div>
  {{ end }}
  {{ if .Remove }}
  <div id="pkgremove" class="listing">
    <h3 title="Removed packages">Removed packages</h3>
    <ul>
      {{ range .Remove }}
      <li><a href="/info?pkg={{ .Name }}&db=local">{{ .Name }}-{{ .Version }}</a></li>
      {{ end }}
    </ul>
  </div>
  {{ end }}
</div>
{{ end }}

{{/* vim: set ts=2 sw=2 et tw=0 ft=gotplhtml: */}}

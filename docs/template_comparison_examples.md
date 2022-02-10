pongo2 with `gets` and pattern
```jinja
{% if getenv("ENV_VARS_SECRET") %}
  {% for kv in gets(printf("/%s/*", getenv("ENV_VARS_SECRET"))) %}
export {{ kv.Key | base }}={{ kv.Value }}
  {% endfor %}
{% else %}
  {% for key in ls("/") %}
export {{ key }}={{ getv(printf("/%s", key)) }}
  {% endfor %}
{% endif %}
```

pongo2 `gets` ( pattern ) VS `ls+getv`
```jinja
{% if getenv("ENV_VARS_SECRET") %}
  {% for kv in gets(printf("/%s/*", getenv("ENV_VARS_SECRET"))) %}
export {{ kv.Key | base }}={{ kv.Value }}
  {% endfor %}
{% endif %}

{% comment %} SAME RESULT AS BELOW {% endcomment %}

{% if getenv("ENV_VARS_SECRET") %}
  {% for key in ls(printf("/%s/", getenv("ENV_VARS_SECRET"))) %}
export {{key}}={{ getv(printf("/%s/%s", getenv("ENV_VARS_SECRET"), key)) }}
  {% endfor %}
{% endif %}

```

go text/template with `gets` and pattern
```go
{{- if getenv "ENV_VARS_SECRET" }}
# secrets from $ENV_VARS_SECRET
{{- range gets ( printf "/%s/*" (getenv "ENV_VARS_SECRET") ) }}
export {{ base .Key }}={{ .Value }}
{{- end }}
{{ else }}
# secret old style, from field value
{{ range $key := ls "/" -}}
export {{ $key }}={{ getv ( printf "/%s" $key ) }}
{{ end -}}
{{ end -}}
```


go text/template `gets` (pattern) VS `ls+get`
```go
{{- if getenv "ENV_VARS_SECRET" }}
# secrets from $ENV_VARS_SECRET
{{- range gets ( printf "/%s/*" (getenv "ENV_VARS_SECRET") ) }}
export {{ base .Key }}={{ .Value }}
{{- end }}
{{- end }}

{{/* SAME AS BELOW */}}

{{- if getenv "ENV_VARS_SECRET" }}
{{- range $key := ls ( printf "/%s/" (getenv "ENV_VARS_SECRET") ) }}
export {{ $key }}={{ getv ( printf "/%s/%s" (getenv "ENV_VARS_SECRET") $key ) }}
{{- end }}
{{- end }}
```

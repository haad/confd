
# Templating text golan / pongo2

There is two option for templating, you can choose from golang text/template extended with sprig ( you know it from helm) and pongo2 which is kind of dialect of jinja2, you can know it from django templates. Both of share some basic functions to fetch data from backands and their differ in manipulation and filters functions.

## Shared data feching functions

		"exists" -  Exists checks for the existence of key in the stor
		"ls"
		"lsdir"
		"get" -  return the KVPair associated with key
		"gets" - returns a KVPair for all nodes with keys matching pattern
		"getv" - gets the value associated with key.
		"getvs" - return list of all values matching pattern

see examples for usage

## Golang specific functions

check [Sprig Docs](http://masterminds.github.io/sprig/)

backward compatible and depricated confd functions
        "json"
        "jsonArray"
        "map"
        "getenv"
        "datetime"
        "toUpper"
        "toLower"
        "lookupIP"
        "lookupIPV4"
        "lookupIPV6"
        "lookupSRV"
        "fileExists"
        "base64Encode"
        "base64Decode"
        "parseBool"
        "reverse"
        "sortByLength"
        "sortKVByLength"
        "seq"
        "printf"
        "unixTS"
        "dateRFC3339"

Beaware there was some naming conflicts between confd original function and sprig, for example split function is now inherited from sprig. Look at [sprig docs](http://masterminds.github.io/sprig/string_slice.html#splitlist-and-split), so split needs to be replaced with splitList sprig function in template to achive old behavior.



https://github.com/flosch/pongo2/blob/master/docs/filters.md

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

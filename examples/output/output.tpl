{{range .}}
##  {{.name}}
{{if .description}}{{.description}}{{end}}
{{range .errors}}
{{if .message}}{{.message}}{{end}}
{{if .location}}{{.location}}{{end}}
{{if .level}}{{.level}}{{end}}
{{if .custom}}{{.custom}}{{end}}
{{end}}
{{end}}

{{range .errors}}
  Name: {{.name}}
  {{if .data_file}}Data: {{.data_file}}{{end}}
{{end}}

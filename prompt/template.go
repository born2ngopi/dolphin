package prompt

import (
	"bytes"
	"text/template"
)

const (
	PROMPT = `i have a function like this

{{.Function}}

{{with .Structs}}
and i have a struct like this

  {{range .}}
type {{.Name}} struct {
  {{range .Fields}}
  {{.Name}} {{.Type}}
  {{end}}
}
  {{end}}
{{end}}

{{with .Mocks}}
and i use mock {{.Name}} and the dir is {{.Dir}}
{{end}}

can you write a unit test for this function?`
)

type Template struct {
	Function string
	Structs  []Struct
	Mock     Mock
}

type Struct struct {
	Name   string
	Fields []StructField
}

type StructField struct {
	Name string
	Type string
}

type Mock struct {
	Name string
	Dir  string
}

func (p Template) Generate() (string, error) {
	tmpl, err := template.New("prompt").Parse(PROMPT)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, p)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

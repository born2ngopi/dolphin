package prompt

import (
	"bytes"
	"go/ast"
	"text/template"
)

const (
	PROMPT = `U are senior and expert golang developer, You are given the task of creating a unit test with a specified format, and set prepare function only if we need mocking some function, you are not allow to mock standard library from golang.
You are asked to create a unit test by testing every possible possibility that occurs according to the parameters and function body given.
You have to write down all the positive and negative cases.
You are also given a list of structs used and list interface, please only use the list of structs that we will provide or structs in the function code, you are not allowed to create new structs outside of those we provide.
And you are only given the opportunity to respond to the unit test code, without information text, import package  or other meaningless, only function unit test code.


function code is
{{.SourceCode}}
{{if .IsMethod}}
and struct for method is

type {{.StuctMethod.Name}} struct {
	{{range .StuctMethod.Fields}}
	{{.Name}} {{.Type}}
	{{end}}
}

{{with .InterfaceMethod}}
and interface on field struct is

{{range .}}
type {{.Name}} interface {
	{{range .Methods}}
	{{.}}
	{{end}}
}
{{end}}

{{end}}

{{end}}

{{with .Structs}}
and have some related structs like this

{{range .}}
type {{.Name}} struct {
	{{range .Fields}}
	{{.Name}} {{.Type}}
	{{end}}
}
from {{.From}}
{{end}}

{{end}}

{{with .Mock}}
and use mock {{.Name}} and the dir is {{.Dir}} and please dont mock standard library
{{end}}

the template for unit test must be like this
func Test_[function_name](t *testing.T) {

	// add some preparation code here include mock, var, and etc

	// add schenario here with []struct
	/*
		example:
		type arg struct {
			// this field must be parameter function
		}

		tests := []struct{
			name string
			arg arg // arg is parameter function,
			wantError error
			wantResponse [response function]
			prepare func([parameter function]) // prepare for expected mock function
		}{
			{
				// fill hire with success scenario and posibility negative/error scenario
			}
		}
	/*

	// looping schenario here and test the function
	/*
		example:
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T){
				// some test logic here

				// check if wantError not nil then check error message is same or not
				// if wantError nil check result function
				// check result and error with assert
			})
		}
	/*
}
	`
)

type Template struct {
	// Function string
	SourceCode      string
	IsMethod        bool
	StuctMethod     Struct
	InterfaceMethod []Interface
	Structs         []Struct
	Mock            Mock
}

type Struct struct {
	Name   string
	From   string
	Fields []StructField
}

type Interface struct {
	Name    string
	Methods []string
}

type StructField struct {
	Name    string
	Type    string
	TypeExp ast.Expr
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

package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/born2ngopi/dolphin/prompt"
)

// variable Struct is a list of struct on the code
var Struct = make(map[string]prompt.Struct)
var Interface = make(map[string]prompt.Interface)

func prepareStruct(dir string) error {

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		// check if is not file and not .go extention
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		packageName := file.Name.Name

		for _, decl := range file.Decls {

			// check type struct declaration
			// example if we have struct like this
			// type User struct {
			//   Name string
			//   Age int
			// }
			// then store it to Struct map
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							structName := typeSpec.Name.Name
							_struct := prompt.Struct{
								Name: structName,
							}
							for _, field := range structType.Fields.List {
								for _, name := range field.Names {
									fieldName := name.Name
									fieldType := getTypeString(field.Type)
									_struct.Fields = append(_struct.Fields, prompt.StructField{
										Name:    fieldName,
										Type:    fieldType,
										TypeExp: field.Type,
									})
								}
							}
							Struct[packageName+structName] = _struct
						} else if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							interfaceName := typeSpec.Name.Name
							_interface := prompt.Interface{
								Name: interfaceName,
							}

							for _, method := range interfaceType.Methods.List {
								methodName := method.Names[0].Name

								// check parameter on method
								// _type, ok := method.Type.(*ast.FuncType)

								methodType := getTypeString(method.Type)

								m := strings.ReplaceAll(methodType, "func", fmt.Sprintf("func %s", methodName))
								// fmt.Println("======> ", m)
								_interface.Methods = append(_interface.Methods, m)
							}

							Interface[packageName+interfaceName] = _interface
						}
					}
				}
			}
		}

		return nil
	})

	return err

}

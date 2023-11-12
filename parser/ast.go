package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"

	"github.com/born2ngopi/dolpin/prompt"
)

func readFileToPrompt(path, funcName, modulePath, dir, mockLib, mockDir string) (prompts []prompt.Template, err error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// prepare import path
	importPath := make(map[string]string)

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok == token.IMPORT {
				for _, spec := range genDecl.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						importPath[importSpec.Name.Name] = importSpec.Path.Value
					}
				}
			}
		}
	}

	// check function
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {

			if funcName != "" && funcName != funcDecl.Name.Name {
				continue
			}

			var _prompt prompt.Template

			// read function code
			sourceCode := getSourceCode(funcDecl.Pos(), funcDecl.End(), fset)
			// read code on this line
			if sourceCode == "" {
				continue
			}
			_prompt.Function = sourceCode

			body := funcDecl.Body

			for _, stmt := range body.List {
				if declStmt, ok := stmt.(*ast.DeclStmt); ok {

					_structs, err := getStructFromStatement(declStmt, importPath)
					if err != nil {
						return nil, err
					}
					if _structs != nil {
						_prompt.Structs = append(_prompt.Structs, _structs...)
					}
				}
			}

			_prompt.Mock = prompt.Mock{
				Name: mockLib,
				Dir:  mockDir,
			}

			prompts = append(prompts, _prompt)

		}
	}

	return prompts, nil
}

// getStructFromStatement is a function to get struct declaration from statement
// example:
// var user = user.User{}
// and User have field
//   - Name string
//   - Age int
//
// then this function will return
//
//	[]prompt.Struct{
//	   {
//	      Name: "User",
//	      Fields: []prompt.Field{
//	         {
//	            Name: "Name",
//	            Type: "string",
//	         },
//	         {
//	            Name: "Age",
//	            Type: "int",
//	         },
//	      },
//	   },
func getStructFromStatement(decl ast.Decl, importPath map[string]string) ([]prompt.Struct, error) {
	var structs []prompt.Struct

	if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				for range valueSpec.Names {
					// check if type variable is struct
					if valueSpec.Type != nil {
						if selExp, ok := valueSpec.Type.(*ast.SelectorExpr); ok {
							// check if struct is from import
							if ident, ok := selExp.X.(*ast.Ident); ok {
								if importPath[ident.Name] != "" {

									err := filepath.Walk(importPath[ident.Name], func(path string, info os.FileInfo, _ error) error {
										// check is file with .go extension
										if info.IsDir() || filepath.Ext(path) != ".go" {
											return nil
										}

										importFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
										if err != nil {
											return err
										}

										structName := selExp.Sel.Name

										if _struct, ok := Struct[structName]; ok {
											structs = append(structs, _struct)
										} else {

											structFieldMap := findStructFields(importFile, structName)

											_struct := prompt.Struct{
												Name: structName,
											}
											if len(structFieldMap) > 0 {
												for fieldName, fieldType := range structFieldMap {
													_struct.Fields = append(_struct.Fields, prompt.StructField{
														Name: fieldName,
														Type: fieldType,
													})
												}
											}

											Struct[structName] = _struct
											structs = append(structs, _struct)
										}

										return nil
									})
									if err != nil {
										return nil, err
									}

								}
							} else {
								// read field struct
								// TODO: check if struct is from same file or sampe package
							}
						}
					}
				}
			}
		}
	}

	return structs, nil
}

func findStructFields(file *ast.File, structName string) map[string]string {
	structFieldMap := make(map[string]string)

	// Inspeksi semua deklarasi
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if typeSpec.Name.Name == structName {
						// Jika ditemukan struct, inspeksi field-nya
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							for _, field := range structType.Fields.List {
								for _, name := range field.Names {
									fieldName := name.Name
									fieldType := getTypeString(field.Type)
									structFieldMap[fieldName] = fieldType
								}
							}
						}
					}
				}
			}
		}
	}

	return structFieldMap
}

func getTypeString(expr ast.Expr) string {
	return types.ExprString(expr)
}

func getSourceCode(start, end token.Pos, fset *token.FileSet) string {
	startOffset := fset.Position(start).Offset
	endOffset := fset.Position(end).Offset

	file, err := os.ReadFile(fset.Position(start).Filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(file)[startOffset:endOffset]
}

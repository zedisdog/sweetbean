/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/spf13/cobra"
	"github.com/zedisdog/sweetbean/errx"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// seederCmd represents the seeder command
var seederCmd = &cobra.Command{
	Use:   "seeder",
	Short: "create seeder",
	Long:  `create seeder.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &packages.Config{
			Mode:  packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps | packages.NeedModule,
			Dir:   cmd.Flag("dir").Value.String(),
			Tests: false,
		}
		pkg, err := packages.Load(c)
		if err != nil {
			panic(err)
		}
		Obj := pkg[0].Types.Scope().Lookup(cmd.Flag("type").Value.String())
		if Obj == nil {
			panic(errx.New("type is not exists"))
		}
		s := Obj.Type().Underlying().(*types.Struct)
		for i := 0; i < s.NumFields(); i++ {
			if _, ok := s.Field(i).Type().(*types.Basic); !ok {
				continue
			}
			// println(s.Field(i).Name() + ":" + s.Field(i).Type().String())
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", "package seed", 0)
		if err != nil {
			panic(err)
		}
		astutil.AddImport(fset, f, Obj.Pkg().Path())
		astutil.AddImport(fset, f, "gorm.io/gorm")

		f.Decls = append(f.Decls, &ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%sSeed", Obj.Name())),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("db"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("gorm"),
									Sel: ast.NewIdent("DB"),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("model"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent("entity." + Obj.Name()),
							},
						},
					},
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("db"),
								Sel: ast.NewIdent("FirstOrCreate"),
							},
							Args: []ast.Expr{
								&ast.UnaryExpr{
									X:  ast.NewIdent("model"),
									Op: token.AND,
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
						},
					},
				},
			},
		})

		var output []byte
		buffer := bytes.NewBuffer(output)
		err = format.Node(buffer, fset, f)
		if err != nil {
			panic(err)
		}
		println(buffer.String())
		// b, err := pkg[0].MarshalJSON()
		// if err != nil {
		// 	panic(err)
		// }
		// println(string(b))
	},
}

func init() {
	createCmd.AddCommand(seederCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seederCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seederCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	seederCmd.Flags().StringP("dir", "d", ".", "module dir path")
	seederCmd.Flags().StringP("type", "t", "", "type for generate")
	seederCmd.MarkFlagRequired("type")
}

type Visitor struct {
	ObjName string
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		n.Decls = append(n.Decls, &ast.FuncDecl{
			Name: ast.NewIdent(fmt.Sprintf("%sSeed", v.ObjName)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("db"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("gorm"),
									Sel: ast.NewIdent("DB"),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("model"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent("entity." + v.ObjName),
							},
						},
					},
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("db"),
								Sel: ast.NewIdent("FirstOrCreate"),
							},
							Args: []ast.Expr{
								&ast.UnaryExpr{
									X:  ast.NewIdent("model"),
									Op: token.AND,
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
						},
					},
				},
			},
		})
	}
	return v
}

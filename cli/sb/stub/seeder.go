package stub

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/zedisdog/sweetbean/errx"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// MakeSeeder generate seed code by given struct
//
//	dir: the directory of package in which the struct, that you want to generate seeder code for.
//	typ: the sturct you want to generate seeder code for.
func MakeSeeder(dir string, typ string) (code []byte, err error) {
	c := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps | packages.NeedModule,
		Dir:   dir,
		Tests: false,
	}

	pkgs, err := packages.Load(c)
	if err != nil {
		err = errx.Wrap(err, "load pkg error")
		return
	}
	Obj := pkgs[0].Types.Scope().Lookup(typ)
	if Obj == nil {
		err = errx.New("type is not exists")
		return
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package seed", 0)
	if err != nil {
		err = errx.Wrap(err, "parse file error")
		return
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
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.SelectorExpr{
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
							Sel: ast.NewIdent("Error"),
						},
					},
				},
			},
		},
	})

	var tmp []byte
	buffer := bytes.NewBuffer(tmp)
	err = format.Node(buffer, fset, f)
	if err != nil {
		err = errx.Wrap(err, "format error")
	}

	code = buffer.Bytes()
	// s := Obj.Type().Underlying().(*types.Struct)
	// for i := 0; i < s.NumFields(); i++ {
	// 	if _, ok := s.Field(i).Type().(*types.Basic); !ok {
	// 		continue
	// 	}
	// 	// println(s.Field(i).Name() + ":" + s.Field(i).Type().String())
	// }

	return
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zedisdog/sweetbean/cli/sb/stub"
	"github.com/zedisdog/sweetbean/cli/sb/utils"
	"github.com/zedisdog/sweetbean/errx"
)

// seederCmd represents the seeder command
var seederCmd = &cobra.Command{
	Use:   "seeder",
	Short: "create seeder",
	Long:  `create seeder.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, modulePath, entityPath, err := utils.ParseDir(cmd)
		if err != nil {
			panic(errx.Wrap(err, "parse dir error"))
		}

		code, err := stub.MakeSeeder(entityPath, cmd.Flag("type").Value.String())
		if err != nil {
			panic(errx.Wrap(err, "generate code error"))
		}

		fpath := fmt.Sprintf("%s/%s.%s", fmt.Sprintf("%s/infra/database/seed", modulePath), strings.ToLower(cmd.Flag("type").Value.String()), ".go")
		err = utils.CreateFile(fpath, code)
		if err != nil {
			panic(errx.Wrap(err, "write file error"))
		}
	},
}

func init() {
	generateCmd.AddCommand(seederCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seederCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seederCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// seederCmd.Flags().StringP("dir", "d", ".", "module dir path")
	seederCmd.Flags().StringP("module", "m", "", "module")
	seederCmd.Flags().StringP("type", "t", "", "type for generate")
	seederCmd.MarkFlagRequired("type")
	seederCmd.MarkFlagRequired("module")
}

// type Visitor struct {
// 	ObjName string
// }

// func (v *Visitor) Visit(node ast.Node) ast.Visitor {
// 	switch n := node.(type) {
// 	case *ast.File:
// 		n.Decls = append(n.Decls, &ast.FuncDecl{
// 			Name: ast.NewIdent(fmt.Sprintf("%sSeed", v.ObjName)),
// 			Type: &ast.FuncType{
// 				Params: &ast.FieldList{
// 					List: []*ast.Field{
// 						{
// 							Names: []*ast.Ident{
// 								ast.NewIdent("db"),
// 							},
// 							Type: &ast.StarExpr{
// 								X: &ast.SelectorExpr{
// 									X:   ast.NewIdent("gorm"),
// 									Sel: ast.NewIdent("DB"),
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Results: &ast.FieldList{
// 					List: []*ast.Field{
// 						{
// 							Type: ast.NewIdent("error"),
// 						},
// 					},
// 				},
// 			},
// 			Body: &ast.BlockStmt{
// 				List: []ast.Stmt{
// 					&ast.AssignStmt{
// 						Lhs: []ast.Expr{
// 							ast.NewIdent("model"),
// 						},
// 						Tok: token.DEFINE,
// 						Rhs: []ast.Expr{
// 							&ast.CompositeLit{
// 								Type: ast.NewIdent("entity." + v.ObjName),
// 							},
// 						},
// 					},
// 					&ast.ExprStmt{
// 						X: &ast.CallExpr{
// 							Fun: &ast.SelectorExpr{
// 								X:   ast.NewIdent("db"),
// 								Sel: ast.NewIdent("FirstOrCreate"),
// 							},
// 							Args: []ast.Expr{
// 								&ast.UnaryExpr{
// 									X:  ast.NewIdent("model"),
// 									Op: token.AND,
// 								},
// 							},
// 						},
// 					},
// 					&ast.ReturnStmt{
// 						Results: []ast.Expr{
// 							ast.NewIdent("nil"),
// 						},
// 					},
// 				},
// 			},
// 		})
// 	}
// 	return v
// }

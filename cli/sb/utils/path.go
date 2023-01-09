package utils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zedisdog/sweetbean/errx"
	"golang.org/x/tools/go/packages"
)

func ParseDir(cmd *cobra.Command) (root string, modulePath string, entityPath string, err error) {
	c := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps | packages.NeedModule,
		Dir:   ".",
		Tests: false,
	}

	pkgs, err := packages.Load(c)
	if err != nil {
		errx.Wrap(err, "load package error")
		return
	}
	root = pkgs[0].Module.Dir
	modulePath = fmt.Sprintf("%s/internal/module/%s", root, cmd.Flag("module").Value.String())
	entityPath = fmt.Sprintf("%s/domain/entity", modulePath)

	return
}

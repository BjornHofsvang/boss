package initialize

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/hashload/boss/internal/pkg/environment"

	"github.com/hashload/boss/internal/pkg/bosspackage"

	"github.com/hashload/boss/internal/pkg/input"
	"github.com/spf13/cobra"
)

type InitConfig struct {
	env   *environment.Env
	quiet bool
}

func NewInitializeCommand(env *environment.Env) *cobra.Command {
	cmdConfig := InitConfig{
		env: env,
	}
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Long:  `This command initialize a new project`,
		Run: func(cmd *cobra.Command, args []string) {
			cmdConfig.runInitialize()
			cmdConfig.queryInput()
		},
	}
}

func (i *InitConfig) runInitialize() {
	i.printHead()

}

func (i *InitConfig) printHead() {
	fmt.Printf(`
This utility will walk you through creating a boss.json file.
It only covers the most common items, and tries to guess sensible defaults.
		 
Use '%s install <pkg>' afterwards to install a package and
save it as a dependency in the boss.json file.
Press ^C at any time to quit.%s`, i.env.App, "\n\n\n")
}

func (i *InitConfig) queryInput() {
	bossPkg := bosspackage.LoadOrNew(i.env.BossFilePath())
	var folderName = ""
	rxp, err := regexp.Compile(`^.+\` + string(filepath.Separator) + `([^\\]+)$`)
	if err == nil {
		allString := rxp.FindAllStringSubmatch(i.env.WorkDir(), -1)
		folderName = allString[0][1]
	}

	bossPkg.Name = folderName
	bossPkg.Version = "1.0.0"
	bossPkg.MainSrc = "./"

	if !i.quiet {
		bossPkg.Name = input.GetTextOrDef("package name", bossPkg.Name)
		bossPkg.Homepage = input.GetTextOrDef("homepage", bossPkg.Homepage)
		bossPkg.Version = input.GetTextOrDef("version", bossPkg.Version)
		bossPkg.Description = input.GetTextOrDef("description", bossPkg.Description)
		bossPkg.MainSrc = input.GetTextOrDef("source folder", bossPkg.MainSrc)

	}

	bossPkg.Save()
}

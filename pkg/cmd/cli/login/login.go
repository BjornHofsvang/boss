package login

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/hashload/boss/internal/pkg/environment"
	"github.com/tcnksm/go-input"

	"github.com/hashload/boss/internal/pkg/configuration"
	"github.com/hashload/boss/internal/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewLoginCommand(env *environment.Env) *cobra.Command {
	var removeLogin bool
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	var cmd = &cobra.Command{
		Use:   "login",
		Short: "Register login to repo",
		Run: func(cmd *cobra.Command, args []string) {
			if removeLogin {
				delete(env.Config.Auth, args[0])
				env.Config.SaveConfiguration()
				return
			}
			var auth *configuration.Auth
			var repo string
			if len(args) > 0 && args[0] != "" {
				repo = args[0]
				auth = env.Config.Auth[args[0]]
			} else {
				response, err := ui.Ask("Url to login (ex: github.com)", &input.Options{
					Required: true,
					Loop:     true,
				})
				utils.CheckError(err)
				auth = env.Config.Auth[response]
			}

			if auth == nil {
				auth = &configuration.Auth{}
			}

			response, err := ui.Select("Use SSH", []string{"yes", "no"}, &input.Options{
				Default:  "yes",
				Loop:     true,
				Required: true,
			})
			utils.CheckError(err)
			auth.UseSsh = response == "yes"
			if auth.UseSsh {
				sshPath, err := ui.Ask("Path of ssh private key", &input.Options{
					Default:  getSSHKeyPath(),
					Loop:     true,
					Required: true,
				})
				utils.CheckError(err)
				auth.Path = sshPath
			} else {
				user, _ := ui.Ask("Username", &input.Options{
					Required: true,
					Loop:     true,
				})
				pass, _ := ui.Ask("Password", &input.Options{
					Mask:     true,
					Required: true,
					Loop:     true,
				})

				auth.SetUser(user)
				auth.SetPassword(pass)
			}
			env.Config.Auth[repo] = auth
			env.Config.SaveConfiguration()

		},
	}
	cmd.Flags().BoolVarP(&removeLogin, "rm", "", false, "remove login")

	return cmd
}

func getSSHKeyPath() string {
	usr, e := user.Current()
	if e != nil {
		log.Fatal(e)
	}
	return filepath.Join(usr.HomeDir, ".ssh", "id_rsa")
}

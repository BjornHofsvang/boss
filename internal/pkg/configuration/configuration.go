package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitSsh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hashload/boss/internal/pkg/crypto"
	"github.com/hashload/boss/internal/pkg/machineinfo"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	GlobalMode          bool             `json:"-"`
	DebugMode           bool             `json:"-"`
	Internal            bool             `json:"-"`
	Id                  string           `json:"id"`
	Auth                map[string]*Auth `json:"auth"`
	PurgeTime           time.Duration    `json:"purge_after"`
	InternalRefreshRate time.Duration    `json:"internal_refresh_rate"`
	LastPurge           time.Time        `json:"last_purge_cache"`
	LastInternalUpdate  time.Time        `json:"last_internal_update"`
	DelphiPath          string           `json:"delphi_path,omitempty"`
	ConfigVersion       int64            `json:"config_version"`
	GitEmbedded         bool             `json:"git_embedded"`
}

type Auth struct {
	UseSsh bool   `json:"usesSsh"`
	Path   string `json:"path,omitempty"`
	User   string `json:"user,omitempty"`
	Pass   string `json:"pass,omitempty"`
}

func LoadConfiguration() *Configuration {
	c := &Configuration{
		DebugMode:           false,
		GlobalMode:          false,
		PurgeTime:           3 * 24 * time.Hour,
		InternalRefreshRate: 5 * 24 * time.Hour,
		Auth:                make(map[string]*Auth),
		Id:                  machineinfo.Md5MachineID,
		GitEmbedded:         true,
	}

	buff, err := ioutil.ReadFile(filepath.Join(c.filePath(), "boss.cfg.json"))
	if err != nil {
		log.Errorf("Fail to load cfg %s", err)
		c.SaveConfiguration()
		return c
	}

	err = json.Unmarshal(buff, c)
	if err != nil {
		log.Fatalf("Fail to load cfg %s", err)
		return c
	}

	if c.Id != machineinfo.Md5MachineID {
		log.Error("Failed to load auth... recreate login accounts")
		c.Id = machineinfo.Md5MachineID
		c.Auth = make(map[string]*Auth)
		c.SaveConfiguration()
	}
	c.SaveConfiguration()
	return c
}

func (c *Configuration) filePath() string {
	homeDir := os.Getenv("BOSS_HOME")
	if homeDir == "" {
		systemHome, err := homedir.Dir()
		homeDir = systemHome
		if err != nil {
			log.Fatal("Error to get cache paths", err)
		}

		homeDir = filepath.FromSlash(homeDir)
	}

	return filepath.Join(homeDir, ".boss")
}
func (c *Configuration) BindFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&c.GlobalMode, "global", "g", false, "Global mode")
	cmd.PersistentFlags().BoolVarP(&c.DebugMode, "debug", "d", false, "Debug mode")
}

func (c *Configuration) AuthMethod(repo string) transport.AuthMethod {
	auth := c.Auth[repo]
	if auth == nil {
		return nil
	} else if auth.UseSsh {
		pem, e := ioutil.ReadFile(auth.Path)
		if e != nil {
			log.Fatalf("Fail to open ssh key %s", e)
		}
		var signer ssh.Signer

		signer, e = ssh.ParsePrivateKey(pem)

		if e != nil {
			panic(e)
		}
		return &gitSsh.PublicKeys{User: "git", Signer: signer}

	} else {
		return &http.BasicAuth{Username: auth.GetUser(), Password: auth.GetPassword()}
	}
}

func (c *Configuration) SaveConfiguration() {
	jsonString, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		log.Fatal("Error on parse config file", err.Error())
	}

	err = os.MkdirAll(c.filePath(), 0755)
	if err != nil {
		log.Fatal("Failed on create path", c.filePath(), err.Error())
	}

	configPath := filepath.Join(c.filePath(), "boss.cfg.json")
	f, err := os.Create(configPath)
	if err != nil {
		log.Fatal("Failed on create file ", configPath, err.Error())
		return
	}

	defer f.Close()

	_, err = f.Write(jsonString)
	if err != nil {
		log.Fatal("Failed on write cache file", err.Error())
	}
}

func (a *Auth) GetUser() string {
	if ret, err := crypto.Decrypt(machineinfo.MachineIDByte, a.User); err != nil {
		log.Error("Fail to decrypt user.")
		return ""
	} else {
		return ret
	}
}

func (a *Auth) GetPassword() string {
	if ret, err := crypto.Decrypt(machineinfo.MachineIDByte, a.Pass); err != nil {
		log.Error("Fail to decrypt pass.", err)
		return ""
	} else {
		return ret
	}
}

func (a *Auth) SetUser(user string) {
	if cUSer, err := crypto.Encrypt(machineinfo.MachineIDByte, user); err != nil {
		log.Error("Fail to crypt user.", err)
	} else {
		a.User = cUSer
	}
}

func (a *Auth) SetPassword(pass string) {
	if cPass, err := crypto.Encrypt(machineinfo.MachineIDByte, pass); err != nil {
		log.Error("Fail to crypt user.", err)
	} else {
		a.Pass = cPass
	}
}

package environment

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cfg "github.com/hashload/boss/internal/pkg/configuration"
	"github.com/hashload/boss/internal/pkg/utils"
	"github.com/mitchellh/go-homedir"
)

type Env struct {
	Config *cfg.Configuration
	App    string
}

func InitializeEnv(config *cfg.Configuration, app string) *Env {
	return &Env{
		Config: config,
		App:    app,
	}
}

func (e *Env) BossHome() string {

	homeDir := os.Getenv("BOSS_HOME")

	if homeDir == "" {
		systemHome, e := homedir.Dir()
		homeDir = systemHome
		if e != nil {
			utils.CheckError(fmt.Errorf("Error to get cache paths: %s", e.Error()))
		}

		homeDir = filepath.FromSlash(homeDir)
	}
	return filepath.Join(homeDir, DotBoss)
}

func (e *Env) HashedDelphiPath() string {
	hasher := md5.New()

	_, err := hasher.Write([]byte(strings.ToLower(e.Config.DelphiPath)))
	utils.CheckError(err)

	hashString := hex.EncodeToString(hasher.Sum(nil))
	if e.Config.Internal {
		hashString = BossInternalDir + hashString
	}
	return hashString
}

func (e *Env) WorkDir() string {
	if e.Config.GlobalMode {
		return filepath.Join(e.BossHome(), FolderDependencies, e.HashedDelphiPath())
	} else {
		dir, err := os.Getwd()
		utils.CheckError(err)
		return dir
	}
}

func (e *Env) GetCacheDir() string {
	return filepath.Join(e.BossHome(), "cache")
}

func (e *Env) BossFilePath() string {
	return filepath.Join(e.WorkDir(), FilePackage)
}

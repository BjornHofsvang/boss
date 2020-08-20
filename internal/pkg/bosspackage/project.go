package bosspackage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hashload/boss/internal/pkg/utils"
)

type BossPackage struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Homepage     string            `json:"homepage"`
	MainSrc      string            `json:"mainsrc"`
	Projects     []string          `json:"projects"`
	Scripts      map[string]string `json:"scripts,omitempty"`
	Dependencies map[string]string `json:"dependencies"`
	Locked       *ProjectLock      `json:"-"`
	path         string            `json:"-"`
}

func MakeNew() *BossPackage {
	project := BossPackage{
		Dependencies: make(map[string]string),
		Projects:     []string{},
		Locked:       &ProjectLock{},
	}

	return &project
}

func LoadOrNew(path string) *BossPackage {
	project := MakeNew()
	project.path = path
	project.Locked = LoadProjectLock(project)

	return project
}

func (p *BossPackage) Save() {
	m, err := json.MarshalIndent(p, "", "\t")
	utils.CheckError(err)

	err = ioutil.WriteFile(p.path, m, os.ModePerm)
	utils.CheckError(err)
}

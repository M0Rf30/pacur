package suse

import (
	"path/filepath"

	"github.com/M0Rf30/pacur/constants"
	"github.com/M0Rf30/pacur/utils"
)

type SuseProject struct {
	Name       string
	Root       string
	MirrorRoot string
	BuildRoot  string
	Path       string
	Distro     string
	Release    string
}

func (p *SuseProject) getBuildDir() (path string, err error) {
	path = filepath.Join(p.BuildRoot, p.Distro+"-"+p.Release)

	err = utils.MkdirAll(path)
	if err != nil {
		return
	}

	return
}

func (p *SuseProject) Prep() (err error) {
	buildDir, err := p.getBuildDir()
	if err != nil {
		return
	}

	keyPath := filepath.Join(p.Path, "..", "sign.key")
	exists, err := utils.Exists(keyPath)
	if err != nil {
		return
	}

	if exists {
		err = utils.CopyFile("", keyPath, buildDir, true)
		if err != nil {
			return
		}
	}

	err = utils.RsyncExt(p.Path, buildDir, ".rpm")
	if err != nil {
		return
	}

	return
}

func (p *SuseProject) Create() (err error) {
	buildDir, err := p.getBuildDir()
	if err != nil {
		return
	}

	err = utils.Exec("", "docker", "run", "--rm", "-t", "-v",
		buildDir+":/pacur:Z", constants.DockerOrg+p.Distro+"-"+p.Release,
		"create", p.Distro+"-"+p.Release, p.Name)
	if err != nil {
		return
	}

	err = utils.Rsync(filepath.Join(buildDir, "zypper"),
		filepath.Join(p.MirrorRoot, "zypper"))
	if err != nil {
		return
	}

	return
}

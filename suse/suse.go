package suse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dropbox/godropbox/container/set"
	"github.com/M0Rf30/pacur/constants"
	"github.com/M0Rf30/pacur/pack"
	"github.com/M0Rf30/pacur/utils"
)

type Suse struct {
	Pack         *pack.Pack
	suseDir      string
	buildDir     string
	buildRootDir string
	rpmsDir      string
	sourcesDir   string
	specsDir     string
	srpmsDir     string
}

func (s *Suse) getRpmPath() (path string, err error) {
	archs, err := ioutil.ReadDir(s.rpmsDir)
	if err != nil {
		return
	}

	if len(archs) < 1 {
		return
	}
	archPath := filepath.Join(s.rpmsDir, archs[0].Name())

	rpms, err := ioutil.ReadDir(archPath)
	if err != nil {
		return
	}

	if len(rpms) < 1 {
		return
	}
	path = filepath.Join(archPath, rpms[0].Name())

	return
}

func (s *Suse) getDepends() (err error) {
	if len(s.Pack.MakeDepends) == 0 {
		return
	}

	args := []string{
		"-n",
		"install",
		"-y",
	}
	args = append(args, s.Pack.MakeDepends...)

	err = utils.Exec("", "zypper", args...)
	if err != nil {
		return
	}

	return
}

func (s *Suse) getFiles() (files []string, err error) {
	backup := set.NewSet()
	paths := set.NewSet()

	for _, path := range s.Pack.Backup {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		backup.Add(path)
	}

	rpmPath, err := s.getRpmPath()
	if err != nil {
		return
	}

	output, err := utils.ExecOutput(s.Pack.PackageDir, "rpm", "-qlp", rpmPath)
	if err != nil {
		return
	}

	for _, path := range strings.Split(output, "\n") {
		if len(path) < 1 || strings.Contains(path, ".build-id") {
			continue
		}

		paths.Remove(filepath.Dir(path))
		paths.Add(path)
	}

	for pathInf := range paths.Iter() {
		path := pathInf.(string)

		if backup.Contains(path) {
			path = `%config "` + path + `"`
		} else {
			path = `"` + path + `"`
		}

		files = append(files, path)
	}

	return
}

func (s *Suse) createSpec(files []string) (err error) {
	path := filepath.Join(s.specsDir, s.Pack.PkgName+".spec")

	release := "%{?dist}"
	if s.Pack.Distro == "amazonlinux" && s.Pack.Release == "1" {
		release = ".amzn1"
	} else if s.Pack.Distro == "amazonlinux" && s.Pack.Release == "2" {
		release = ".amzn2"
	} else if s.Pack.Distro == "centos" && s.Pack.Release == "7" {
		release = ".el7.centos"
	} else if s.Pack.Distro == "oraclelinux" && s.Pack.Release == "7" {
		release = ".el7.oraclelinux"
	}

	data := ""
	data += fmt.Sprintf("Name: %s\n", s.Pack.PkgName)
	data += fmt.Sprintf("Summary: %s\n", s.Pack.PkgDesc)
	data += fmt.Sprintf("Version: %s\n", s.Pack.PkgVer)
	data += fmt.Sprintf("Release: %s%s\n", s.Pack.PkgRel, release)
	data += fmt.Sprintf("Group: %s\n", ConvertSection(s.Pack.Section))
	data += fmt.Sprintf("URL: %s\n", s.Pack.Url)
	data += fmt.Sprintf("License: %s\n", s.Pack.License)
	data += fmt.Sprintf("Packager: %s\n", s.Pack.Maintainer)

	for _, pkg := range s.Pack.Provides {
		data += fmt.Sprintf("Provides: %s\n", pkg)
	}

	for _, pkg := range s.Pack.Conflicts {
		data += fmt.Sprintf("Conflicts: %s\n", pkg)
	}

	for _, pkg := range s.Pack.Depends {
		data += fmt.Sprintf("Requires: %s\n", pkg)
	}

	for _, pkg := range s.Pack.MakeDepends {
		data += fmt.Sprintf("BuildRequires: %s\n", pkg)
	}

	data += "\n"

	if len(s.Pack.PkgDescLong) > 0 {
		data += "%description\n"
		for _, line := range s.Pack.PkgDescLong {
			data += line + "\n"
		}
		data += "\n"
	}

	data += "%install\n"
	data += fmt.Sprintf("rsync -a -A %s/ $RPM_BUILD_ROOT/\n",
		s.Pack.PackageDir)
	data += "\n"

	data += "%files\n"
	if len(files) == 0 {
		data += "/\n"
	} else {
		for _, line := range files {
			data += line + "\n"
		}
	}
	data += "\n"

	if len(s.Pack.PreInst) > 0 {
		data += "%pre\n"
		for _, line := range s.Pack.PreInst {
			data += line + "\n"
		}
		data += "\n"
	}

	if len(s.Pack.PostInst) > 0 {
		data += "%post\n"
		for _, line := range s.Pack.PostInst {
			data += line + "\n"
		}
		data += "\n"
	}

	if len(s.Pack.PreRm) > 0 {
		data += "%preun\n"
		data += "if [[ \"$1\" -ne 0 ]]; then exit 0; fi\n"
		for _, line := range s.Pack.PreRm {
			data += line + "\n"
		}
		data += "\n"
	}

	if len(s.Pack.PostRm) > 0 {
		data += "%postun\n"
		data += "if [[ \"$1\" -ne 0 ]]; then exit 0; fi\n"
		for _, line := range s.Pack.PostRm {
			data += line + "\n"
		}
	}

	err = utils.CreateWrite(path, data)
	if err != nil {
		return
	}

	fmt.Println(data)

	return
}

func (s *Suse) rpmBuild() (err error) {
	err = utils.Exec(s.specsDir, "rpmbuild", "--define",
		"_topdir "+s.suseDir, "-ba", s.Pack.PkgName+".spec")
	if err != nil {
		return
	}

	return
}

func (s *Suse) Prep() (err error) {
	err = s.getDepends()
	if err != nil {
		return
	}

	return
}

func (s *Suse) makeDirs() (err error) {
	s.suseDir = filepath.Join(s.Pack.Root, "suse")
	s.buildDir = filepath.Join(s.suseDir, "BUILD")
	s.buildRootDir = filepath.Join(s.suseDir, "BUILDROOT")
	s.rpmsDir = filepath.Join(s.suseDir, "RPMS")
	s.sourcesDir = filepath.Join(s.suseDir, "SOURCES")
	s.specsDir = filepath.Join(s.suseDir, "SPECS")
	s.srpmsDir = filepath.Join(s.suseDir, "SRPMS")

	for _, path := range []string{
		s.suseDir,
		s.buildDir,
		s.buildRootDir,
		s.rpmsDir,
		s.sourcesDir,
		s.specsDir,
		s.srpmsDir,
	} {
		err = utils.ExistsMakeDir(path)
		if err != nil {
			return
		}
	}

	return
}

func (s *Suse) clean() (err error) {
	pkgPaths, err := utils.FindExt(s.Pack.Home, ".rpm")
	if err != nil {
		return
	}

	match, ok := constants.ReleasesMatch[s.Pack.FullRelease]
	if !ok {
		return
	}

	for _, pkgPath := range pkgPaths {
		if strings.Contains(filepath.Base(pkgPath), match) {
			_ = utils.Remove(pkgPath)
		}
	}

	return
}

func (s *Suse) copy() (err error) {
	archs, err := ioutil.ReadDir(s.rpmsDir)
	if err != nil {
		return
	}

	for _, arch := range archs {
		err = utils.CopyFiles(filepath.Join(s.rpmsDir, arch.Name()),
			s.Pack.Home, false)
		if err != nil {
			return
		}
	}

	return
}

func (s *Suse) remDirs() {
	os.RemoveAll(s.suseDir)
}

func (s *Suse) Build() (err error) {
	err = s.makeDirs()
	if err != nil {
		return
	}
	defer s.remDirs()

	err = s.createSpec([]string{})
	if err != nil {
		return
	}

	err = s.rpmBuild()
	if err != nil {
		return
	}

	files, err := s.getFiles()
	if err != nil {
		return
	}

	if len(files) == 0 {
		return
	}

	s.remDirs()
	err = s.makeDirs()
	if err != nil {
		return
	}

	err = s.createSpec(files)
	if err != nil {
		return
	}

	err = s.rpmBuild()
	if err != nil {
		return
	}

	err = s.clean()
	if err != nil {
		return
	}

	err = s.copy()
	if err != nil {
		return
	}

	return
}

package constants

import (
	"strings"

	"github.com/dropbox/godropbox/container/set"
)

const (
	DockerOrg = "m0rf30/pacur-"
)

var (
	Releases = [...]string{
		"archlinux",
		"amazonlinux-1",
		"amazonlinux-2",
		"fedora-30",
		"fedora-31",
		"centos-7",
		"debian-jessie",
		"debian-stretch",
		"debian-buster",
		"opensuse",
		"oraclelinux-7",
		"ubuntu-xenial",
		"ubuntu-bionic",
		"ubuntu-eoan",
	}
	ReleasesMatch = map[string]string{
		"archlinux":      "",
		"amazonlinux-1":  ".amzn1.",
		"amazonlinux-2":  ".amzn2.",
		"fedora-29":      ".fc29.",
		"fedora-30":      ".fc30.",
		"centos-7":       ".el7.centos.",
		"debian-jessie":  ".jessie_",
		"debian-stretch": ".stretch_",
		"debian-buster":  ".buster_",
		"opensuse":       ".opensuse_",
		"oraclelinux-7":  ".el7.oraclelinux.",
		"ubuntu-trusty":  ".trusty_",
		"ubuntu-xenial":  ".xenial_",
		"ubuntu-bionic":  ".bionic_",
		"ubuntu-disco":   ".disco_",
		"ubuntu-eoan":    ".eoan_",
	}
	DistroPack = map[string]string{
		"archlinux":   "pacman",
		"amazonlinux": "redhat",
		"fedora":      "redhat",
		"centos":      "redhat",
		"debian":      "debian",
		"opensuse":    "suse",
		"oraclelinux": "redhat",
		"ubuntu":      "debian",
	}
	Packagers = [...]string{
		"apt",
		"pacman",
		"yum",
		"zypper",
	}

	ReleasesSet    = set.NewSet()
	Distros        = []string{}
	DistrosSet     = set.NewSet()
	DistroPackager = map[string]string{}
	PackagersSet   = set.NewSet()
)

func init() {
	for _, release := range Releases {
		ReleasesSet.Add(release)
		distro := strings.Split(release, "-")[0]
		Distros = append(Distros, distro)
		DistrosSet.Add(distro)
	}

	for _, distro := range Distros {
		packager := ""

		switch DistroPack[distro] {
		case "debian":
			packager = "apt"
		case "pacman":
			packager = "pacman"
		case "redhat":
			packager = "yum"
		case "suse":
			packager = "zypper"
		default:
			panic("Failed to find packager for distro")
		}

		DistroPackager[distro] = packager
	}

	for _, packager := range Packagers {
		PackagersSet.Add(packager)
	}
}

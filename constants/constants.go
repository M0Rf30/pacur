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
		"amazonlinux-1",
		"amazonlinux-2",
		"archlinux",
		"astralinux",
		"fedora-30",
		"fedora-31",
		"centos-7",
		"debian-buster",
		"debian-stretch",
		"opensuse",
		"oraclelinux-7",
		"ubuntu-bionic",
		"ubuntu-eoan",
		"ubuntu-xenial",
	}
	ReleasesMatch = map[string]string{
		"amazonlinux-1":  ".amzn1.",
		"amazonlinux-2":  ".amzn2.",
		"archlinux":      "",
		"astralinux":     ".astra2_",
		"fedora-30":      ".fc30.",
		"fedora-31":      ".fc31.",
		"centos-7":       ".el7.centos.",
		"debian-stretch": ".stretch_",
		"debian-buster":  ".buster_",
		"opensuse":       ".opensuse_",
		"oraclelinux-7":  ".el7.oraclelinux.",
		"ubuntu-bionic":  ".bionic_",
		"ubuntu-eoan":    ".eoan_",
		"ubuntu-xenial":  ".xenial_",
	}
	DistroPack = map[string]string{
		"amazonlinux": "redhat",
		"archlinux":   "pacman",
		"astralinux":  "debian",
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

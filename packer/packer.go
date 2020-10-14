package packer

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/M0Rf30/pacur/constants"
	"github.com/M0Rf30/pacur/debian"
	"github.com/M0Rf30/pacur/pack"
	"github.com/M0Rf30/pacur/pacman"
	"github.com/M0Rf30/pacur/redhat"
	"github.com/M0Rf30/pacur/suse"
)

type Packer interface {
	Prep() error
	Build() error
}

func GetPacker(pac *pack.Pack, distro, release string) (
	pcker Packer, err error) {

	switch constants.DistroPack[distro] {
	case "pacman":
		pcker = &pacman.Pacman{
			Pack: pac,
		}
	case "debian":
		pcker = &debian.Debian{
			Pack: pac,
		}
	case "redhat":
		pcker = &redhat.Redhat{
			Pack: pac,
		}
	case "suse":
		pcker = &suse.Suse{
			Pack: pac,
		}
	default:
		system := distro
		if release != "" {
			system += "-" + release
		}

		err = &UnknownSystem{
			errors.Newf("packer: Unkown system %s", system),
		}
		return
	}

	return
}

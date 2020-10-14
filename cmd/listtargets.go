package cmd

import (
	"fmt"
	"github.com/M0Rf30/pacur/constants"
)

func ListTargets() (_ error) {
	for _, release := range constants.Releases {
		fmt.Println(release)
	}

	return
}

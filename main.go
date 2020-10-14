package main

import (
	"github.com/M0Rf30/pacur/cmd"
)

func main() {
	err := cmd.Parse()
	if err != nil {
		panic(err)
	}
}

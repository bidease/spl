package main

import (
	"log"

	"github.com/docopt/docopt-go"
)

type config struct {
	Conf         string
	Hostid       string
	Delete       bool
	S, Sshkeys   bool
	L, Locations bool
	C, Cloud     bool
	D, Dedicated bool
}

var conf config
var usage = `
Usage:
	spl --help
	spl (c | cloud) [<hostid>] [--delete] [--conf FILE]
	spl (d | dedicated) [<hostid>] [--conf FILE]
	spl (s | sshkeys) [--conf FILE]
	spl (l | locations) [--conf FILE]

Options:
	--help       show this help message and exit
	--conf FILE  path to config file [default: ~/.spl.yml]
`

func readArgs() {
	docopts, _ := docopt.ParseDoc(usage)
	if err := docopts.Bind(&conf); err != nil {
		log.Fatal(err)
	}
}

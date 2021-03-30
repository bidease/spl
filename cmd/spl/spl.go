package main

import (
	"log"

	"github.com/bidease/spl"
)

func main() {
	switch {
	case conf.Sshkeys || conf.S:
		printSSHKeys()
	case conf.Locations || conf.L:
		printLocations()
	case (conf.Dedicated || conf.D) && conf.Hostid == "":
		printDedicatedServers()
	case (conf.Dedicated || conf.D) && conf.Hostid != "":
		getDedicatedServersDescribe()
	case (conf.Cloud || conf.C) && conf.Hostid == "":
		printCloudInstances()
	case (conf.Cloud || conf.C) && conf.Hostid != "":
		getCloudInstanceDescribe()
	}
}

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	readArgs()
	spl.Conf.Read(conf.Conf)
}

package main

import (
	"log"
	"os"

	"github.com/bidease/spl"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	app := cli.NewApp()
	app.HideHelp = true
	app.Author = "Konstantin Kruglov"
	app.Email = "kruglovk@gmail.com"
	app.Version = "2.0.0"
	app.Before = initial
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: "~/.spl.yml",
			Usage: "path to config file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "sshkeys",
			Aliases: []string{"s"},
			Action:  printSSHKeys,
			Usage:   "print SSH keys",
		},
		{
			Name:    "locations",
			Aliases: []string{"l"},
			Action:  printLocations,
			Usage:   "print locations",
		},
		{
			Name:    "cloud",
			Usage:   "cloud instances",
			Aliases: []string{"c"},
			Action:  printCloudInstances,
			Subcommands: []cli.Command{
				{
					Name:    "describe",
					Aliases: []string{"d"},
					Action:  getCloudInstanceDescribe,
					Usage:   "print cloud instances",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "instance ID",
						},
					},
				},
			},
		},
		{
			Name:    "dedicated",
			Usage:   "print dedicated servers",
			Aliases: []string{"d"},
			Action:  printDedicatedServes,
			Subcommands: []cli.Command{
				{
					Name:    "describe",
					Aliases: []string{"d"},
					Usage:   "print describe",
					Action:  getDedicatedServersDescribe,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "server ID",
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func initial(c *cli.Context) error {
	spl.Conf.Read(c.GlobalString("conf"))
	return nil
}

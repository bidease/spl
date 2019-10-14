package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bidease/spl/cloud"
	"github.com/bidease/spl/common"
	"github.com/bidease/spl/config"
	"github.com/bidease/spl/tools"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	app := cli.NewApp()
	app.HideHelp = true
	app.Author = "Konstantin Kruglov"
	app.Email = "kruglovk@gmail.com"
	app.Version = "1.3.0"
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
			Usage:   "porint SSH keys",
			Action:  common.PrintSSHKeys,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "public",
					Usage: "show public kyes",
				},
			},
		},
		{
			Name:    "hardware",
			Aliases: []string{"h"},
			Usage:   "list/detail",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Action:  listHosts,
					Usage:   "print exists hardware servers",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "price",
							Usage: "show price",
						},
						cli.BoolFlag{
							Name:  "cpu",
							Usage: "show CPU name",
						},
					},
				},
				{
					Name:    "info",
					Aliases: []string{"i"},
					Action:  detailInfoAboutServer,
					Usage:   "detail infomation about server",
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name:  "id",
							Usage: "ID of server",
						},
					},
				},
			},
		},
		{
			Name:    "cloud",
			Aliases: []string{"c"},
			Usage:   "info",
			Subcommands: []cli.Command{
				{
					Name:    "regions",
					Aliases: []string{"r"},
					Action:  cloud.PrintRegions,
					Usage:   "print available regions",
				},
				{
					Name:    "images",
					Aliases: []string{"i"},
					Action:  cloud.PrintImages,
					Usage:   "print available images in region for orders",
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:  "id",
							Usage: "ID of region",
							Value: 0,
						},
					},
				},
				{
					Name:    "servers",
					Aliases: []string{"s"},
					Action:  cloud.PrintCloudServers,
					Usage:   "print available configs cloud servers in region for orders",
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:  "id",
							Usage: "ID of region",
							Value: 0,
						},
					},
				},
				{
					Name:    "create",
					Aliases: []string{"c"},
					Action:  cloud.CreateCloudServer,
					Usage:   "create cloud server",
					Flags: []cli.Flag{
						cli.UintFlag{
							Name:  "regionID",
							Usage: "ID of region",
						},
						cli.StringFlag{
							Name:  "imageID",
							Usage: "ID of image",
						},
						cli.StringFlag{
							Name:  "configID",
							Usage: "ID of config (VCPUs, RAM, etc)",
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "host name",
						},
						cli.StringFlag{
							Name:  "fingerprint",
							Usage: "use an user name",
						},
						cli.BoolFlag{
							Name:  "password",
							Usage: "use password",
						},
						cli.UintFlag{
							Name:  "backups",
							Usage: "number last backups",
							Value: 5,
						},
						cli.BoolFlag{
							Name:  "gpn",
							Usage: "global private network",
						},
					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Action:  cloud.PrintExistsCloudServers,
					Usage:   "print exists cloud servers",
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Action:  cloud.DeleteInstanse,
					Usage:   "delete exists cloud server",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "ID of deleting cloud server",
						},
						cli.StringFlag{
							Name:  "token",
							Usage: "one-time token",
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
	config.Options.Read(c.GlobalString("conf"))
	config.Options.Check()
	return nil
}

func listHosts(c *cli.Context) {
	var hs hosts
	tools.GetRequest(hostsURL, &hs)
	table := tablewriter.NewWriter(os.Stdout)
	lField := []string{"id", "host name", "location", "public ip", "private ip"}

	if c.Bool("cpu") {
		lField = append(lField, "cpu")
	}

	if c.Bool("price") {
		lField = append(lField, "price")
	}

	table.SetHeader(lField)
	var price float64
	for _, vol := range hs.Data {
		row := []string{
			fmt.Sprint(vol.ID),
			vol.Title,
			shortLocation(vol.Location.Name),
			getIP("public", &vol.commonInfoHost),
			getIP("private", &vol.commonInfoHost),
		}

		if c.Bool("cpu") {
			var h host
			tools.GetRequest(fmt.Sprintf(hostURL, vol.ID), &h)
			row = append(row, h.Data.Server.ChassisModelCPUName)
		}

		if c.Bool("price") {
			curPrice := getPrice(vol.ID)
			price = price + curPrice
			row = append(row, fmt.Sprintf("%.2f", getPrice(vol.ID)))
		}
		table.Append(row)
	}
	table.Render()

	var b balance
	tools.GetRequest(balanceURL, &b)
	fmt.Printf("Total servers: %d\n", hs.NumFound)
	fmt.Printf("The total cost of servers: %f\n", price)
	fmt.Printf("Balance: %s\n", b.Data.Balance)
	fmt.Printf("Estimated balance: %s\n\n", b.Data.EstimatedBalance)
}

func detailInfoAboutServer(c *cli.Context) {
	serverID := c.Uint64("id")
	var h host
	tools.GetRequest(fmt.Sprintf(hostURL, serverID), &h)
	s := getServices(serverID)

	fmt.Printf("Price: %.2f\n", getPrice(h.Data.ID))
	fmt.Println("Start rent:", getStartRentHost(&s))
	if h.Data.ScheduledReleaseAt != "" {
		fmt.Printf("Scheduled release at: %s\n", h.Data.ScheduledReleaseAt)
	}
	fmt.Printf("Host name: %s\n", h.Data.commonInfoHost.Title)
	fmt.Printf("Location: %s\n", h.Data.Location.Name)
	fmt.Printf("Timezone: %s\n", h.Data.Location.Timezone)
	fmt.Printf("OS: %s %s %s\n", h.Data.OS.Name, h.Data.OS.Version, h.Data.OS.Arch)
	if h.Data.OSReinstallation {
		fmt.Println("OS reinstallation: yes")
	}
	if h.Data.HasDRAC {
		fmt.Printf("Has DRAC: %s\n", map[bool]string{true: "yes", false: "no"}[h.Data.HasDRAC])
		fmt.Printf("DRAC is enable: %s\n", map[string]string{"enabled": "yes", "disabled": "no"}[h.Data.DRACIsEnable])
	}
	fmt.Printf("CPU: %s\n", h.Data.Server.ChassisModelCPUName)
	fmt.Printf("Configuration: %s\n", h.Data.Server.Configuration)
	fmt.Printf("Networks:\n")

	for _, v := range h.Data.commonInfoHost.Networks {
		fmt.Printf("  IP: %s\n", v.HostIP)
		fmt.Printf("  Netmask: %s\n", v.Netmask)
		fmt.Printf("  Type: %s\n", v.PoolType)
	}

	fmt.Println("Uplinks:")
	for _, v := range s.Data {
		if v.Type == 11 {
			fmt.Printf("  %s (%.2f %s)\n", v.Description, v.Price, v.Currency)
		}
	}
}

func shortLocation(l string) string {
	if len(l) == 0 {
		return "UNKNOWN"
	}

	return strings.Split(l, " ")[0]
}

func getIP(t string, n *commonInfoHost) string {
	for _, v := range n.Networks {
		if t == v.PoolType {
			return v.HostIP
		}
	}
	return "UNKNOWN"
}

func getPrice(id uint64) (price float64) {
	for _, v := range getServices(id).Data {
		price = price + v.Price
	}
	return
}

func getServices(id uint64) (s services) {
	tools.GetRequest(fmt.Sprintf(servicesURL, id), &s)
	return
}

func getStartRentHost(s *services) string {
	for _, service := range s.Data {
		if service.Type == 1 {
			return service.DateStart
		}
	}
	return "UNKNOWN"
}

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	app.Version = "1.1.0"
	app.Before = initial
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "~/.spl.yml",
			Usage: "path to config file",
		},
	}
	app.Commands = []cli.Command{
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
	tools.Request(hostsURL, &hs)
	table := tablewriter.NewWriter(os.Stdout)
	lField := []string{"id", "host name", "location", "public ip", "private ip"}

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
			getIP("private", &vol.commonInfoHost)}

		if c.Bool("price") {
			curPrice := getPrice(vol.ID)
			price = price + curPrice
			row = append(row, fmt.Sprint(getPrice(vol.ID)))
		}
		table.Append(row)
	}
	table.Render()

	var b balance
	tools.Request(balanceURL, &b)
	fmt.Printf("Total servers: %d\n", hs.NumFound)
	fmt.Printf("The total cost of servers: %f\n", price)
	fmt.Printf("Balance: %s\n", b.Data.Balance)
	fmt.Printf("Estimated balance: %s\n\n", b.Data.EstimatedBalance)
}

func detailInfoAboutServer(c *cli.Context) {
	serverID := c.Int64("id")
	var h host
	tools.Request(fmt.Sprintf(hostURL, serverID), &h)

	fmt.Println()
	fmt.Printf("Price: %.2f\n", getPrice(h.Data.ID))
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
		fmt.Println()
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
	var s services
	tools.Request(fmt.Sprintf(servicesURL, id), &s)

	for _, v := range s.Data {
		price = price + v.Price
	}
	return
}

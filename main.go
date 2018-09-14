package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	yml "gopkg.in/yaml.v2"
)

var conf config

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	app := cli.NewApp()
	app.HideHelp = true
	app.Author = "Konstantin Kruglov"
	app.Email = "kruglovk@gmail.com"
	app.Version = "1.0.0"
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
				},
				{
					Name:    "detail",
					Aliases: []string{"d"},
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
	conf.read(c.GlobalString("conf"))
	conf.check()
	return nil
}

func listHosts(c *cli.Context) {
	var hs hosts
	request(hostsURL, &hs)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "host name", "location", "public ip", "private ip"})

	for _, vol := range hs.Data {
		table.Append([]string{
			fmt.Sprint(vol.ID),
			vol.Title,
			shortLocation(vol.Location.Name),
			getIP("public", &vol.commonInfoHost),
			getIP("private", &vol.commonInfoHost)})
	}
	table.Render()

	var b balance
	request(balanceURL, &b)
	fmt.Printf("Total servers: %d\nBalance: %s\nEstimated balance: %s\n\n", hs.NumFound, b.Data.Balance, b.Data.EstimatedBalance)
}

func detailInfoAboutServer(c *cli.Context) {
	serverID := c.Int64("id")
	var h host
	request(fmt.Sprintf(hostURL, serverID), &h)

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

func request(path string, out interface{}) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://portal.servers.com/rest/%s", path), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Email", conf.Email)
	req.Header.Set("X-User-Token", conf.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	err = json.Unmarshal(body, &out)
	if err != nil {
		log.Fatalln(err)
	}
}

type config struct {
	Email string
	Token string
}

func (c *config) read(f string) {
	if !path.IsAbs(f) && f[:1] == "~" {
		f = path.Join(os.Getenv("HOME"), f[1:])
	}

	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("Read file %s is failed: %s", f, err)
	}

	err = yml.Unmarshal(bytes, c)
	if err != nil {
		log.Fatalf("Read config is failed: %s", err)
	}
}

func (c *config) check() {
	if c.Email == "" {
		log.Fatalln("Email not defined")
	}
	if c.Token == "" {
		log.Fatalln("Token not defined")
	}
}

func shortLocation(l string) string {
	switch l {
	case "AMS1 (Amsterdam Metropolitan Area, The Netherlands)":
		return "EU/Amsterdam"
	case "LUX2 (Roost, Luxembourg)":
		return "EU/Roost"
	case "MOW1 (Moscow, Russian Federation)":
		return "EU/Moscow"
	case "DFW1 (Dallas–Fort Worth, TX, USA)":
		return "USA/Dallas"
	case "DFW2 (Dallas–Fort Worth, TX, USA)":
		return "USA/Dallas"
	default:
		return "UNKNOWN"
	}
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
	request(fmt.Sprintf(servicesURL, id), &s)

	for _, v := range s.Data {
		price = price + v.Price
	}
	return
}

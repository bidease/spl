package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/olekukonko/tablewriter"
	yml "gopkg.in/yaml.v2"
)

var (
	confFile  = flag.String("c", "~/.spl.yml", "path to config file")
	sLocation = flag.Bool("s", true, "print short location name")
	serverID  = flag.Int64("S", 0, "detail information about the server")
	conf      config
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	flag.Parse()

	conf.read(*confFile)
	conf.check()

	if *serverID > 0 {
		detailInfoAboutServer()
		return
	}

	listHosts()
}

func listHosts() {
	var hs hosts
	request(hostsURL, &hs)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "HostName", "Location", "Public IP", "Private IP"})

	for _, vol := range hs.Data {
		if *sLocation {
			vol.Location.Name = shortLocation(vol.Location.Name)
		}
		table.Append([]string{
			fmt.Sprint(vol.ID),
			vol.Title,
			vol.Location.Name,
			getIP("public", &vol.commonInfoHost),
			getIP("private", &vol.commonInfoHost)})
	}

	// table.SetBorder(false)
	table.Render()

	var b balance
	request(balanceURL, &b)
	fmt.Printf("Total servers: %d\nBalance: %s\nEstimated balance: %s\n\n", hs.NumFound, b.Data.Balance, b.Data.EstimatedBalance)
}

func detailInfoAboutServer() {
	var h host
	request(fmt.Sprintf(hostURL, *serverID), &h)

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

package cloud

import (
	"fmt"
	"log"
	"os"

	"github.com/bidease/spl/tools"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// PrintExistsCloudServers ..
func PrintExistsCloudServers(c *cli.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	lField := []string{
		"id", "name", "status", "flavor name", "region", "local ip", "internal ip", "external ip",
	}
	table.SetHeader(lField)

	for _, instanse := range getExistsCloudServers() {
		row := []string{
			instanse.ID,
			instanse.Name,
			instanse.Status,
			instanse.FlavorName,
			instanse.RegionName,
			instanse.LocalIP,
			instanse.InternalIP,
			instanse.ExternalIP,
		}
		table.Append(row)
	}
	table.Render()
}

// DeleteInstanse ..
func DeleteInstanse(c *cli.Context) {
	var response struct {
		Success bool
	}

	for _, instanse := range getExistsCloudServers() {
		if instanse.ID == c.String("id") {
			if len(c.String("token")) > 0 {
				token := struct {
					Token string `json:"token"`
				}{
					Token: c.String("token"),
				}
				tools.DeleteRequest(fmt.Sprintf("cloud_computing/regions/%s/instances/%s", instanse.RegionID, c.String("id")), &response, token)
			} else {
				tools.DeleteRequest(fmt.Sprintf("cloud_computing/regions/%s/instances/%s", instanse.RegionID, c.String("id")), &response, nil)
			}
			break
		}
	}

	if !response.Success {
		log.Fatalf("%v\n", response.Success)
	}
}

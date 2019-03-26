package cloud

import (
	"os"

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
			instanse.Region,
			instanse.LocalIP,
			instanse.InternalIP,
			instanse.ExternalIP,
		}
		table.Append(row)
	}
	table.Render()
}

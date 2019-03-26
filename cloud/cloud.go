package cloud

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/urfave/cli"
)

// PrintRegions ..
func PrintRegions(c *cli.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	lField := []string{"id", "name"}
	table.SetHeader(lField)

	for _, item := range getRegions() {
		row := []string{item.ID, item.Name}
		table.Append(row)
	}
	table.Render()
}

// PrintImages ..
func PrintImages(c *cli.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	lField := []string{
		"name", "id", "created at", "display priority", "requires ssh key",
		"min disk", "is windows", "allowed servers",
	}
	table.SetHeader(lField)

	for _, item := range getImages(c.Uint("id")) {
		row := []string{
			item.Name,
			item.ID,
			item.CreatedAt,
			fmt.Sprintf("%d", item.DisplayPriority),
			fmt.Sprintf("%v", item.RequiresSSHKey),
			fmt.Sprintf("%d", item.MinDisk),
			fmt.Sprintf("%v", item.IsWindows),
			fmt.Sprintf("%s", strings.Join(item.AllowedFlavors, ", ")),
		}
		table.Append(row)
	}
	table.Render()
}

// PrintCloudServers ..
func PrintCloudServers(c *cli.Context) {
	table := tablewriter.NewWriter(os.Stdout)
	lField := []string{
		"id", "name", "vcpus", "ram", "disk", "free traffic (over)", "price",
	}
	table.SetHeader(lField)

	for _, item := range getCloudServers(c.Uint("id")) {
		row := []string{
			item.ID,
			item.Name,
			fmt.Sprintf("%d", item.VCPUs),
			fmt.Sprintf("%d MB", item.RAM),
			fmt.Sprintf("%d GB SSD", item.Disk),
			fmt.Sprintf("%d GB (%.2f/GB)", item.FreeTrafficAmountGb, item.OvercommitTrafficPricePerGb),
			fmt.Sprintf("%.2f %s", item.MonthlyPricesPerUnit.Full.Hosting.Total, item.MonthlyPricesPerUnit.Currency),
		}
		table.Append(row)
	}
	table.Render()
}

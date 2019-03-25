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
		"min disk", "is windows", "allowed flavors",
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

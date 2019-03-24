package cloud

import (
	"os"

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

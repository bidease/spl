package common

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// PrintSSHKeys ..
func PrintSSHKeys(c *cli.Context) {
	table := tablewriter.NewWriter(os.Stdout)

	var lField []string

	if c.Bool("public") {
		lField = []string{"name", "fingerprint", "public key"}
	} else {
		lField = []string{"name", "fingerprint"}
	}

	table.SetHeader(lField)

	for _, item := range GetSSHKeys() {
		var row []string
		if c.Bool("public") {
			row = []string{item.Name, item.Fingerprint, item.PublicKey}
		} else {
			row = []string{item.Name, item.Fingerprint}
		}
		table.Append(row)
	}
	table.Render()
}

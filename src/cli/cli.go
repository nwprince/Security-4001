package cli

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/segmentio/go-prompt"

	"nodes"
)

var opts = []string{
	"List",
	"Send cmd to device",
	"Download file from device",
	"Upload file to device",
}

// Shell will create an interactive console
func Shell() {
	for {
		i := prompt.Choose("Choose an option:", opts)
		fmt.Println("picked: ", opts[i])

		switch opts[i] {
		case "List":
			listNodes()
		}
	}
}

func listNodes() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"GUID", "PLATFORM", "ARCH", "USER", "HOST", "STATUS"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	for k, v := range nodes.Nodes {
		table.Append([]string{k.String(), v.Platform, v.Architecture, v.UserName, v.HostName, nodes.GetStatus(k)})
	}
	fmt.Println()
	table.Render()
	fmt.Println()
}

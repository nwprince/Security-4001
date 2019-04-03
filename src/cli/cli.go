package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/go-prompt"

	"nodes"
)

var mainOpts = []string{
	"List",
	"Interact",
}

var interactOpts = []string{
	"Send cmd to device",
	"Download file from device",
	"Upload file to device",
}

// Shell will create an interactive console
func Shell() {
	for {
		i := prompt.Choose("Choose an option:", mainOpts)

		switch i {
		case 0:
			listNodes()
		case 1:
			interact()
		}
	}
}

func listNodes() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"GUID", "PLATFORM", "ARCH", "USER", "HOST", "STATUS"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	var counts int
	for k, v := range nodes.Nodes {
		counts = 0
		table.Append([]string{k.String(), v.Platform, v.Architecture, v.UserName, v.HostName, nodes.GetStatus(k)})
		counts++
	}
	fmt.Println()
	table.Render()
	fmt.Println()
}

func readInput(inputMsg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(inputMsg)
	text, _ := reader.ReadString('\n')
	return text
}

func interact() {
	listNodes()
	i := prompt.Choose("Choose an option:", interactOpts)

	switch i {
	case 0:
		sendCmd()
	}

}

func selectNode() uuid.UUID {
	var opts []uuid.UUID
	var stringOpts []string

	for k := range nodes.Nodes {
		opts = append(opts, k)
		stringOpts = append(stringOpts, k.String())
	}
	i := prompt.Choose("Select a device to interact with: ", stringOpts)
	return opts[i]
}

func sendCmd() {
	node := selectNode()
	cmd := readInput("Please enter your command: ")
	nodes.AddJob(node, "cmdString", []string{cmd})

}

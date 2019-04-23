package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		break
	case 1:
		downloadFile()
		break
	case 2:
		uploadFile()
		break
	}

}

func selectNode(blast bool) uuid.UUID {
	var opts []uuid.UUID
	var stringOpts []string

	for k := range nodes.Nodes {
		opts = append(opts, k)
		stringOpts = append(stringOpts, k.String())
	}
	if blast {
		u, _ := uuid.FromString("ffffffff-ffff-ffff-ffff-ffffffffffff")
		opts = append(opts, u)
		stringOpts = append(stringOpts, "ffffffff-ffff-ffff-ffff-ffffffffffff")
	}
	i := prompt.Choose("Select a device to interact with: ", stringOpts)
	return opts[i]
}

func sendCmd() {
	node := selectNode(true)
	cmd := readInput("Please enter your command (absolute path): ")
	nodes.AddJob(node, "cmdString", strings.Fields(cmd))

}

func downloadFile() {
	node := selectNode(true)
	filename := readInput("Please enter the file name (absolute path): ")
	if len(filename) > 0 {
		nodes.AddJob(node, "download", strings.Fields(filename))
	}
}

func uploadFile() {
	node := selectNode(true)
	filename := readInput("Please enter the file name (absolute path): ")
	remotepath := readInput("Remote File Path (absolute path): ")
	args := []string{strings.Fields(filename)[0], strings.Fields(remotepath)[0]}
	if len(filename) > 0 {
		nodes.AddJob(node, "upload", args)
	}

}

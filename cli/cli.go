package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nwprince/Security-4001/nodes"
	"github.com/olekukonko/tablewriter"
	uuid "github.com/satori/go.uuid"
)

var mainOpts = []string{
	"List",
	"Interact",
}

var interactOpts = []string{
	"Send cmd to device",
	"Download file from device",
	"Upload file to device",
	"Execute Script",
}

func createString(prompt string, args ...interface{}) string {
	var s string
	fmt.Printf(prompt+": ", args...)
	fmt.Scanln(&s)
	return s
}

func choose(prompt string, list []string) int {
	fmt.Println()
	for i, val := range list {
		fmt.Printf("  %d) %s\n", i+1, val)
	}

	fmt.Println()
	i := -1

	for {
		s := createString(prompt)

		// index
		n, err := strconv.Atoi(s)
		if err == nil {
			if n > 0 && n <= len(list) {
				i = n - 1
				break
			} else {
				continue
			}
		}

		// value
		i = indexOf(s, list)
		if i != -1 {
			break
		}
	}

	return i
}

func indexOf(s string, list []string) int {
	for i, val := range list {
		if val == s {
			return i
		}
	}
	return -1
}

// Shell will create an interactive console
func Shell() {

	for {
		i := choose("Choose an option:", mainOpts)

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
	i := choose("Choose an option:", interactOpts)

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
	case 3:
		script()
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
	i := choose("Select a device to interact with: ", stringOpts)
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

func script() {
	node := selectNode(true)
	filename := readInput("Enter the name of the script (absolute path): ")
	if len(filename) > 0 {
		nodes.AddJob(node, "script", strings.Fields(filename))
	}
}

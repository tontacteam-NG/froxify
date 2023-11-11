package binary

import (
	"fmt"
	"os"
	"os/exec"
)

func X8(url string) {
	// Command to run
	command := "/Users/nxczje/Documents/go/froxy/binary/x8Mac"
	args := []string{"-u", url, "-w", "/Users/nxczje/Documents/go/froxy/lists/param.txt"}

	// Create a new command
	cmd := exec.Command(command, args...)

	// Set the output to os.Stdout to see the result in the console
	cmd.Stdout = os.Stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

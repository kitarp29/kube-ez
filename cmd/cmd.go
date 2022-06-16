package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func Command(commands string) string {
	command := strings.Fields(commands)
	fmt.Println(command)
	cmd := exec.Command(command[0], command[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		return (err.Error())
	}
	// Print the output
	return (string(stdout))
}

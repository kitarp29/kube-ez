package cmd

import (
	"os/exec"
)

func Event() string {

	cmd := exec.Command("kubectl", "get", "events", "-o", "json")
	stdout, err := cmd.Output()
	// fmt.Println(reflect.TypeOf(stdout))
	if err != nil {
		return (err.Error())
	}
	// Print the output
	return (string(stdout))
}

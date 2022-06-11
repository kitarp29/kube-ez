package cmd

import (
	"fmt"
	"os/exec"
)

func CreateNamespace() {
	app := "kubectl"

	arg0 := "create"
	arg1 := "namespace"
	arg2 := "pk"
	// arg2 := "\n\tfrom"
	// arg3 := "golang"

	cmd := exec.Command(app, arg0, arg1, arg2)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

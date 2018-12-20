package main

import (
	"fmt"
	"os"

	"github.com/lunarway/dextre/cmd"
)

func main() {
	command, err := cmd.NewCommand("dextre")
	if err != nil {
		fmt.Printf("Prerequisuites failed. Error: %v\n", err)
		os.Exit(1)
	}

	command.Execute()

}

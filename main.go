package main

import (
	"fmt"
	"os"

	"github.com/lunarway/dextre/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "<dev-version>"
	commit  = "<unspecified-commit>"
)

func main() {
	command, err := cmd.NewCommand("dextre")
	if err != nil {
		fmt.Printf("Prerequisuites failed. Error: %v\n", err)
		os.Exit(1)
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dextre",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
	command.AddCommand(versionCmd)
	command.Execute()

}

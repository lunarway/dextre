package cmd

import (
	"github.com/CrowdSurge/banner"
	"github.com/lunarway/dextre/cmd/drain"
	"github.com/lunarway/dextre/cmd/roll"
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// NewCommand RootCmd for Cobra
func NewCommand(name string) (*cobra.Command, error) {

	bannerString := banner.PrintS("dextre")

	c := &cobra.Command{
		Use:   name,
		Short: "Small cli tool to move pods in a safe manner",
		Long:  bannerString,
	}

	var kubeConfig string
	var verbose bool
	flags := c.PersistentFlags()
	flags.StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig file")
	flags.BoolVar(&verbose, "verbose", false, "verbose output")

	kubectl, err := kubernetes.NewClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	c.AddCommand(
		drain.NewCommand(kubectl, &verbose),
		roll.NewCommand(kubectl, &verbose),
	)
	return c, nil
}

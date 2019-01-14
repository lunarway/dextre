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
	flags := c.PersistentFlags()
	flags.StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig file")

	kubectl, err := kubernetes.NewClient(kubeConfig)
	if err != nil {
		return nil, err
	}

	c.AddCommand(
		drain.NewCommand(kubectl),
		roll.NewCommand(kubectl),
	)
	return c, nil
}

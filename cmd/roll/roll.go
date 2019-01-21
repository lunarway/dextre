package roll

import (
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// NewCommand sets up the move command
func NewCommand(kubectl *kubernetes.Client, verbose *bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "roll",
		Short: "",
		Long:  "",
	}
	c.AddCommand(
		podsCommand(kubectl, verbose),
		nodesCommand(kubectl, verbose),
	)

	return c
}

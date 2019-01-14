package roll

import (
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/roll"
	"github.com/spf13/cobra"
)

// NewCommand sets up the move command
func nodesCommand(kubectl *kubernetes.Client) *cobra.Command {
	var role string
	var label string

	c := &cobra.Command{
		Use:   "nodes",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return roll.Nodes(kubectl, role, label)
		},
	}
	c.Flags().StringVar(&role, "role", "", "Role type of the nodes to be rolled")
	c.MarkFlagRequired("role")
	c.Flags().StringVar(&label, "label", "", "label of the nodes to be rolled")

	return c
}

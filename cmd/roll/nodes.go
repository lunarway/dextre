package roll

import (
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/roll"
	"github.com/spf13/cobra"
)

// NewCommand sets up the move command
func nodesCommand(kubectl *kubernetes.Client, verbose *bool) *cobra.Command {
	var instanceGroup string
	var cluster string
	var awsRegion string

	c := &cobra.Command{
		Use:   "nodes",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return roll.Nodes(kubectl, instanceGroup, cluster, awsRegion, *verbose)
		},
	}
	c.Flags().StringVar(&instanceGroup, "kops-instance-group", "", "kops instance group to perfrom the rolling on")
	c.MarkFlagRequired("kops-instance-group")
	c.Flags().StringVar(&label, "label", "", "label of the nodes to be rolled")
	c.Flags().StringVar(&cluster, "cluster", "", "the name of the kops cluster")
	c.Flags().StringVar(&awsRegion, "aws-region", "ue-west-1", "the region to use")
	c.MarkFlagRequired("cluster")

	return c
}

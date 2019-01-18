package roll

import (
	"time"

	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/roll"
	"github.com/spf13/cobra"
)

var (
	label       string
	namespace   string
	gracePeriod time.Duration
)

// NewCommand sets up the move command
func podsCommand(kubectl *kubernetes.Client, verbose *bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "pods",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return roll.Pods(kubectl, label, namespace, gracePeriod, *verbose)
		},
	}
	c.Flags().StringVar(&label, "label", "", "The labels that should be restarted on the form: type=service")
	c.Flags().StringVar(&namespace, "namespace", "", "The namespace to search for pods")
	c.MarkFlagRequired("label")
	c.MarkFlagRequired("namespace")
	c.Flags().DurationVar(&gracePeriod, "grace-period", (10 * time.Second), "pod grace-period")

	return c
}

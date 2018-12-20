package restart

import (
	"time"

	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/restart"
	"github.com/spf13/cobra"
)

var (
	label       string
	namespace   string
	gracePeriod time.Duration
)

// NewCommand sets up the move command
func NewCommand(kubectl *kubernetes.Client) *cobra.Command {
	c := &cobra.Command{
		Use:   "restart",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return restart.Run(kubectl, label, namespace, gracePeriod)
		},
	}
	c.Flags().StringVar(&label, "label", "", "The labels that should be restarted on the form: type=service")
	c.Flags().StringVar(&namespace, "namespace", "", "The namespace to search for pods")
	c.MarkFlagRequired("label")
	c.MarkFlagRequired("namespace")
	c.Flags().DurationVar(&gracePeriod, "grace-period", (10 * time.Second), "pod grace-period")

	return c
}

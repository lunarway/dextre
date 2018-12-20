package restart

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/ui"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Run: executes the drain command
func Run(kubectl *kubernetes.Client, label string, namespace string, gracePeriod time.Duration) error {
	// Output the banner
	ui.PrintBanner("dextre")

	pods, err := kubectl.GetPodsWithLabel(label, namespace)
	if err != nil {
		return err
	}

	ui.PrintPodList(pods.Items, "Pods to be restarted", false)

	fmt.Printf("Are you sure you want to restart the pods? ")

	ok, err := ui.AskForConfirmation()
	if err != nil {
		return err
	}

	if !ok {
		// return nil to exit nicely
		return nil
	}

	// restartPods
	fmt.Println("")
	restartPods(kubectl, pods.Items, gracePeriod)

	fmt.Println("")
	color.Green("[✓] All pods restarted!\n")

	return nil
}

func restartPods(kubectl *kubernetes.Client, pods []v1.Pod, gracePeriod time.Duration) error {
	table := ui.NewTable("[-]", "EVICTED", "NEW", "NODE")
	graceP := int64(gracePeriod.Seconds())

	// Evict regular pods first.
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: &graceP}

	// Evict Regular Pods an Wait for New Pod to be ready
	for _, pod := range pods {
		table.PrepareRow()
		err := kubectl.DeletePod(pod, deleteOptions)
		if err != nil {
			return err
		}
		newPod, err := kubectl.DetermineNewPod(pod)
		if newPod != nil {
			err = kubectl.WaitForPodToBeReady(newPod)
			if err != nil {
				return err
			}
			table.CommitRow("[✓]", pod.Name, newPod.Name, newPod.Spec.NodeName)
		}
		table.DiscardRow()
	}
	return nil
}

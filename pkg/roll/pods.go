package roll

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
func Pods(kubectl *kubernetes.Client, label string, namespace string, gracePeriod time.Duration, verbose bool) error {

	pods, err := kubectl.GetPodsWithLabel(label, namespace)
	if err != nil {
		return err
	}

	ui.PrintPodList(pods.Items, "Pods to be restarted", false, verbose)

	fmt.Printf("Are you sure you want to roll the pods? ")

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
	rollPods(kubectl, pods.Items, gracePeriod, verbose)

	fmt.Println("")
	color.Green("[✓] All pods rolled!\n")

	return nil
}

func rollPods(kubectl *kubernetes.Client, pods []v1.Pod, gracePeriod time.Duration, verbose bool) {
	table := ui.NewTable("[-]", "EVICTED", "NEW POD", "NEW NODE", verbose)
	graceP := int64(gracePeriod.Seconds())

	// Evict regular pods first.
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: &graceP}

	// Evict Regular Pods an Wait for New Pod to be ready
	for _, pod := range pods {
		table.PrepareRow()
		err := kubectl.DeletePod(pod, deleteOptions)
		if err != nil {
			table.CommitRow("[-]", pod.Name, "Delete pod",err.Error())
			table.DiscardRow()
			continue
		}
		newPod, err := kubectl.DetermineNewPod(pod)
		if err != nil {
			table.CommitRow("[-]", pod.Name, "Detemine new pod",err.Error())
			table.DiscardRow()
			continue
		}
		if newPod != nil {
			err = kubectl.WaitForPodToBeReady(newPod)
			if err != nil {
				table.CommitRow("[-]", pod.Name, "Wait for ready pod",err.Error())
			} else{
				table.CommitRow("[✓]", pod.Name, newPod.Name, newPod.Spec.NodeName)
			}
		}
		table.DiscardRow()
	}
	return
}
package drain

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	dextreaws "github.com/lunarway/dextre/pkg/aws"
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/ui"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Run: executes the drain command
func Run(kubectl *kubernetes.Client, nodeName string, gracePeriod time.Duration, skipValidation bool) error {
	// Output the banner
	ui.PrintBanner("dextre")

	// Find the node in the cluster
	node, err := kubectl.GetNode(nodeName)
	if err != nil {
		return err
	}

	// Get the pods running on the given node
	pods, err := kubectl.GetPodsOnNode(nodeName)
	if err != nil {
		return err
	}

	// Separate pods into systemPods and regular pods
	var systemPods, regularPods []v1.Pod

	// Group pods into systemPods and Regular pods
	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			systemPods = append(systemPods, pod)
		} else {
			regularPods = append(regularPods, pod)
		}
	}

	// Print the pods to be evicted; both system and regular
	ui.PrintPodList(systemPods, "System pods to be evicted", false)
	ui.PrintPodList(regularPods, "Regular pods to be evicted", true)

	if skipValidation {
		fmt.Printf("Are you sure you want to evict all pods on the node? ")

		ok, err := ui.AskForConfirmation()
		if err != nil {
			return err
		}

		if !ok {
			// return nil to exit nicely
			return nil
		}
	}

	// Cordon the node for in order to not get more pods scheduled
	_, err = kubectl.CordonNode(node)
	if err != nil {
		return err
	}

	fmt.Println("")
	color.Yellow("Cordon\n")
	fmt.Printf("[✓] %s cordoned\n\n", node.ObjectMeta.Name)

	// put this in UI as a header function
	color.Yellow("Evict Regular pods")
	rollPods(kubectl, regularPods, gracePeriod)

	fmt.Println("")
	color.Yellow("Evict System pods")
	rollPods(kubectl, systemPods, gracePeriod)

	fmt.Println("")
	color.Green("[✓] All pods evicted!\n")

	if skipValidation {
		fmt.Println("")
		fmt.Printf("Do you want to continue and terminate the node? ")
		ok, err := ui.AskForConfirmation()
		if err != nil {
			return err
		}

		// user stopped the flow
		if !ok {
			return nil
		}
	}

	fmt.Println("")
	color.Yellow("Node termination:\n")

	// Create the client
	client, err := dextreaws.NewClient("eu-west-1")

	if err != nil {
		return err
	}

	instanceID, err := client.GetInstanceId(nodeName)
	if err != nil {
		return err
	}

	fmt.Printf("%-25s %s\n", "Private DNS:", nodeName)
	fmt.Printf("%-25s %s\n", "Instance ID:", instanceID)

	err = client.TerminateInstance(instanceID)
	if err != nil {
		return err
	}

	fmt.Println("")
	color.Green("[✓] Node has been terminated!\n")

	return nil
}

func rollPods(kubectl *kubernetes.Client, pods []v1.Pod, gracePeriod time.Duration) error {
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

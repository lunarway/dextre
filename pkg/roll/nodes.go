package roll

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lunarway/dextre/pkg/drain"
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/ui"
	"k8s.io/api/core/v1"
)

//Run: executes the drain command
func Nodes(kubectl *kubernetes.Client, role, label string) error {
	// Output the banner
	ui.PrintBanner("dextre")

	nodes, err := kubectl.ListNodes()
	if err != nil {
		return err
	}

	var rollableNodes []v1.Node
	var discardedNodes []v1.Node
	for _, node := range nodes.Items {
		if node.Labels["kubernetes.io/role"] == role {
			if label == "" {
				rollableNodes = append(rollableNodes, node)
			} else {
				labelSlice := strings.Split(label, "=")
				if node.Labels[labelSlice[0]] == labelSlice[1] {
					rollableNodes = append(rollableNodes, node)
				} else {
					discardedNodes = append(discardedNodes, node)
				}
			}
		} else {
			discardedNodes = append(discardedNodes, node)
		}
	}

	fmt.Printf("Nodes to be rolled: %d, Discarded: %d\n", len(rollableNodes), len(discardedNodes))
	fmt.Println("")

	color.Yellow("List of nodes to be rolled:\n")

	for _, node := range rollableNodes {
		fmt.Printf("> %s\n", node.Name)
	}

	fmt.Println("")
	color.Yellow("PROCESS:\n")

	for _, node := range rollableNodes {

		// Drain the Node
		fmt.Println("")
		drain.Run(kubectl, node.Name, 30, true)

		// Wait for the node to be terminated
		err = kubectl.WaitForNodeToTerminate(node)
		if err != nil {
			return err
		}
		color.Green("[✓] Node has been removed from the cluster\n")

		// Get the new list of nodes
		nodes, err = kubectl.ListNodes()
		if err != nil {
			return err
		}

		// Wait and identify the new node
		newNode, err := kubectl.IdentifyNewNode(nodes.Items, role, label)
		if err != nil {
			return err
		}
		color.Green("[✓] %s has been added to the cluster\n", newNode.Name)

		// Wait for the new node to enter Ready state
		err = kubectl.WaitForNewNodeToBeReady(newNode)
		if err != nil {
			return err
		}
		color.Green("[✓] %s is now READY\n", newNode.Name)

	}

	return nil
}

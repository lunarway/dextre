package roll

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	dextreaws "github.com/lunarway/dextre/pkg/aws"
	"github.com/lunarway/dextre/pkg/drain"
	"github.com/lunarway/dextre/pkg/kubernetes"
	"github.com/lunarway/dextre/pkg/ui"
	v1 "k8s.io/api/core/v1"
)

//Run: executes the drain command
func Nodes(kubectl *kubernetes.Client, instanceGroup, cluster string, awsRegion string, verbose bool) error {

	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	// Get nodes in the kubernetes cluster
	nodes, err := kubectl.ListNodes()
	if err != nil {
		return err
	}

	// Find the nodes to roll by matching label and role
	var rollableNodes []v1.Node
	var discardedNodes []v1.Node
	for _, node := range nodes.Items {
		if node.Labels["kops.k8s.io/instancegroup"] == instanceGroup {
			rollableNodes = append(rollableNodes, node)
		} else {
			discardedNodes = append(discardedNodes, node)
		}
	}

	ui.PrintTitle("OVERVIEW:\n", true)
	ui.Print(fmt.Sprintf("Nodes to be rolled: %d", len(rollableNodes)), true)
	for _, node := range rollableNodes {
		ui.Print(fmt.Sprintf("> %s", node.Name), true)
	}
	ui.Print(fmt.Sprintf("Discarded nodes: %d", len(discardedNodes)), true)

	ui.Print("", true)
	ui.PrintTitle("PROGRESS:\n", true)

	// Create AWS client
	client, err := dextreaws.NewClient("eu-west-1")
	if err != nil {
		return err
	}

	for _, node := range rollableNodes {
		spinner.Start()
		ui.PrintTitle(fmt.Sprintf("[-] %s:", node.Name), true)

		// Get the list of current nodes in the cluster
		nodes, err = kubectl.ListNodes()
		if err != nil {
			return err
		}

		// Identify the correct autoscaling group
		asg, err := client.GetAutoScalingGroup(instanceGroup, cluster)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] AWS Autoscaling group located: %s", asg.AutoScalingGroupName), true)

		// Increment the desired number of instances
		err = client.IncrementCapacity(asg)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] Increasing DesiredCapacity of %s from %d nodes to %d nodes and DefaultCooldown to %d", asg.AutoScalingGroupName, asg.DesiredCapacity, asg.DesiredCapacity+1, 0), true)

		// Wait and identify the new node
		newNode, err := kubectl.IdentifyNewNode(nodes.Items, instanceGroup)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] %s added to the cluster", newNode.Name), true)

		// Wait for the new node to enter Ready state
		err = kubectl.WaitForNewNodeToBeReady(newNode)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] %s is now in state: Ready", newNode.Name), true)

		// Drain the Node
		drain.Run(kubectl, node.Name, 30, true, false, awsRegion, verbose)

		// Get the AWS InstanceId
		instanceID, err := client.GetInstanceId(node.Name)
		if err != nil {
			return err
		}

		// Terminate Decrement DesiredCapacity
		err = client.TerminateInstanceDecrementDesiredCapacity(instanceID)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] %s will now be terminated", node.Name), true)

		// Wait for the node to be terminated
		err = kubectl.WaitForNodeToTerminate(node)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] %s is now removed from the cluster", node.Name), true)

		// RESTORE ASG MAXSIZE
		err = client.RestoreValuesForAutoScalingGroup(asg)
		if err != nil {
			return err
		}
		ui.Print(fmt.Sprintf("[✓] Restored AutoScalingGroup %s defaults. MaxSize=%d and DefaultCooldown=%d", asg.AutoScalingGroupName, asg.MaxSize, asg.DefaultCooldown), true)

		// Interval between node roll
		time.Sleep(5 * time.Second)
		spinner.Stop()
	}

	return nil
}

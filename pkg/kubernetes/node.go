package kubernetes

import (
	"time"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// GetNode finds the node in the cluster
func (c *Client) GetNode(nodeName string) (*v1.Node, error) {
	node, err := c.clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the node")
	}
	return node, nil
}

// GetPodsOnNode returns all pods running on the specific node
func (c *Client) GetPodsOnNode(nodeName string) (*v1.PodList, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{
			"spec.nodeName": nodeName,
			"status.phase":  "Running",
		}).String()})

	if err != nil {
		return nil, errors.Wrap(err, "failed to load pods on the node "+nodeName)
	}
	return pods, nil
}

// CordonNode makes the node unschedulable
func (c *Client) CordonNode(node *v1.Node) (*v1.Node, error) {
	node, err := c.clientset.CoreV1().Nodes().Get(node.ObjectMeta.Name, metav1.GetOptions{})
	node.Spec.Unschedulable = true
	n, err := c.clientset.CoreV1().Nodes().Update(node)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cordon node: "+node.Name)
	}
	return n, nil
}

// GetNode finds the node in the cluster
func (c *Client) ListNodes() (*v1.NodeList, error) {
	nodeList, err := c.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find the nodes")
	}
	return nodeList, nil
}

func (c *Client) WaitForNodeToTerminate(node v1.Node) error {
	watcher, err := c.clientset.CoreV1().Nodes().Watch(metav1.ListOptions{
		FieldSelector: "metadata.name=" + node.Name,
	})
	if err != nil {
		return errors.Wrap(err, "cannot create Pod status listener")
	}

	for {
		e := <-watcher.ResultChan()
		if e.Object == nil {
			return errors.Wrap(err, "cannot read object")
		}
		n, ok := e.Object.(*v1.Node)
		if !ok {
			continue
		}

		if n.Name != node.Name {
			continue
		}

		if e.Type == "DELETED" {
			break
		}
	}
	watcher.Stop()
	return nil
}

func (c *Client) IdentifyNewNode(nodes []v1.Node, instanceGroup string) (*v1.Node, error) {
	watcher, err := c.clientset.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create Pod status listener")
	}

	for {
		e := <-watcher.ResultChan()
		if e.Object == nil {
			return nil, errors.Wrap(err, "cannot read object")
		}
		n, ok := e.Object.(*v1.Node)
		if !ok {
			continue
		}

		if e.Type == "ADDED" {
			if !Contains(nodes, n) {
				// Verify that the new node is of the correct type
				if n.Labels["kops.k8s.io/instancegroup"] == instanceGroup {
					return n, nil
				} else {
					// continue if the new node does match role or labels
					continue
				}
			}
		}
	}
	watcher.Stop()
	return nil, nil
}

func (c *Client) WaitForNewNodeToBeReady(node *v1.Node) error {
	node, err := c.clientset.CoreV1().Nodes().Get(node.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to find the node")
	}

	var ready bool
	for {
		node, err = c.clientset.CoreV1().Nodes().Get(node.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to find the node")
		}
		ready = false
		for _, cond := range node.Status.Conditions {
			if cond.Type == v1.NodeReady && cond.Status == v1.ConditionTrue {
				ready = true
				break
			}
		}
		if ready {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func Contains(nodes []v1.Node, n *v1.Node) bool {
	for _, node := range nodes {
		if n.Name == node.Name {
			return true
		}
	}
	return false
}

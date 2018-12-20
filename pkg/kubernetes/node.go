package kubernetes

import (
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

package kubernetes

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPodsWithLabel return all pods with the specified label
func (c *Client) GetPodsWithLabel(label, namespace string) (*v1.PodList, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: label})

	if err != nil {
		return nil, errors.Wrap(err, "failed to load pods with "+label)
	}
	return pods, nil
}

func (c *Client) DeletePod(pod v1.Pod, deleteOptions *metav1.DeleteOptions) error {
	err := c.clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, deleteOptions)
	if err != nil {
		return errors.Wrap(err, "failed to delete pod: "+pod.Name)
	}
	return nil
}

func (c *Client) DetermineNewPod(oldPod v1.Pod) (*v1.Pod, error) {
	// Find all pods with the same labels:
	pods, err := c.clientset.CoreV1().Pods(oldPod.Namespace).List(metav1.ListOptions{
		LabelSelector: "app=" + oldPod.ObjectMeta.Labels["app"]})
	if len(pods.Items) == 0 {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to list pods with labels: "+"app="+oldPod.ObjectMeta.Labels["app"])
	}

	// Sort Pods to find the latest one
	sort.Slice(pods.Items, func(i, j int) bool {
		return pods.Items[i].CreationTimestamp.Time.After(pods.Items[j].CreationTimestamp.Time)
	})

	newPod := pods.Items[0]

	return &newPod, nil
}

func (c *Client) WaitForPodToBeReady(pod *v1.Pod) error {
	watcher, err := c.clientset.CoreV1().Pods(pod.Namespace).Watch(metav1.ListOptions{
		FieldSelector: "metadata.name=" + pod.Name,
	})
	if err != nil {
		return errors.Wrap(err, "cannot create Pod status listener")
	}

	for {
		e := <-watcher.ResultChan()
		if e.Object == nil {
			return errors.Wrap(err, "cannot read object")
		}
		p, ok := e.Object.(*v1.Pod)
		if !ok {
			continue
		}

		if p.Name != pod.Name {
			continue
		}

		if p.Status.Phase == "Running" {
			for i := 0; i <= 30; i++ {
				p, err = c.clientset.CoreV1().Pods(p.Namespace).Get(pod.Name, metav1.GetOptions{})
				if err != nil {
					return errors.Wrap(err, "error retrieving status for new pod")
				}
				if all(p.Status.ContainerStatuses, func(status v1.ContainerStatus) bool { return status.Ready }) {
					break
				}
				time.Sleep(1 * time.Second)
			}
			break
		}
	}
	watcher.Stop()
	return nil
}

func all(vs []v1.ContainerStatus, f func(v1.ContainerStatus) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

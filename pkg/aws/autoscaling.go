package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type AutoScalingGroupStruct struct {
	AutoScalingGroupARN  string
	AutoScalingGroupName string
	DesiredCapacity      int64
	MinSize              int64
	MaxSize              int64
	DefaultCooldown      int64
}

func (c *Client) TerminateInstanceKeepDesiredCapacity(instanceID string) error {
	auto := autoscaling.New(c.Session)
	input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{
		InstanceId:                     aws.String(instanceID),
		ShouldDecrementDesiredCapacity: aws.Bool(false),
	}

	_, err := auto.TerminateInstanceInAutoScalingGroup(input)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) TerminateInstanceDecrementDesiredCapacity(instanceID string) error {
	auto := autoscaling.New(c.Session)
	input := &autoscaling.TerminateInstanceInAutoScalingGroupInput{
		InstanceId:                     aws.String(instanceID),
		ShouldDecrementDesiredCapacity: aws.Bool(true),
	}

	_, err := auto.TerminateInstanceInAutoScalingGroup(input)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAutoScalingGroup(instanceGroup, cluster string) (AutoScalingGroupStruct, error) {
	auto := autoscaling.New(c.Session)
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	asgs, err := auto.DescribeAutoScalingGroups(input)
	if err != nil {
		return AutoScalingGroupStruct{}, err
	}

	for _, a := range asgs.AutoScalingGroups {
		var clusterMatch, instanceGroupMatch bool
		for _, t := range a.Tags {
			// Is it the correct cluster?
			if (*t.Key == "KubernetesCluster") && (*t.Value == cluster) {
				clusterMatch = true
			}
			// Is the label correct?
			if (*t.Key == "k8s.io/cluster-autoscaler/node-template/label/kops.k8s.io/instancegroup") && (*t.Value == instanceGroup) {
				instanceGroupMatch = true
			}
		}
		if instanceGroupMatch && clusterMatch {
			return AutoScalingGroupStruct{
				AutoScalingGroupARN:  *a.AutoScalingGroupARN,
				AutoScalingGroupName: *a.AutoScalingGroupName,
				DesiredCapacity:      *a.DesiredCapacity,
				MinSize:              *a.MinSize,
				MaxSize:              *a.MaxSize,
				DefaultCooldown:      *a.DefaultCooldown,
			}, nil
		}
	}
	return AutoScalingGroupStruct{}, errors.New("not found")
}

func (c *Client) IncrementCapacity(autoScalingGroup AutoScalingGroupStruct) error {
	auto := autoscaling.New(c.Session)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(autoScalingGroup.AutoScalingGroupName),
		DesiredCapacity:      aws.Int64(autoScalingGroup.DesiredCapacity + 1),
		MaxSize:              aws.Int64(autoScalingGroup.DesiredCapacity + 1),
		DefaultCooldown:      aws.Int64(0),
	}
	_, err := auto.UpdateAutoScalingGroup(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RestoreValuesForAutoScalingGroup(autoScalingGroup AutoScalingGroupStruct) error {
	auto := autoscaling.New(c.Session)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(autoScalingGroup.AutoScalingGroupName),
		MaxSize:              aws.Int64(autoScalingGroup.MaxSize),
		DefaultCooldown:      aws.Int64(autoScalingGroup.DefaultCooldown),
	}
	_, err := auto.UpdateAutoScalingGroup(input)
	if err != nil {
		return err
	}

	return nil
}

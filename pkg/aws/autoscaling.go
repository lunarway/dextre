package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func (c *Client) TerminateInstance(instanceID string) error {
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

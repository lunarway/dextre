package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *Client) GetInstanceId(nodeName string) (string, error) {
	// Initialize ec2 service
	svc := ec2.New(c.Session)

	// Specify the instance we need more information about
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(nodeName)},
			},
		},
	}

	// Describe the instance
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		return "", err
	}

	// Retrieve the instance Id
	var instanceID string
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceID = *instance.InstanceId
			break
		}
		break
	}
	return instanceID, nil
}

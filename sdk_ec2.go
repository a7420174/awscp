package awscp

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

func GetReservations(cfg aws.Config, name string, tagKey string) []types.Reservation {
	client := ec2.NewFromConfig(cfg)

	var filterName, filterTag types.Filter
	if name != "" {
		tag1 := "tag:Name"
		filterName = types.Filter{
			Name:   &tag1,
			Values: []string{name},
		}
	}

	if tagKey != "" {
		tag2 := "tag-key"
		filterTag = types.Filter{
			Name:   &tag2,
			Values: []string{tagKey},
		}
	}

	output, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{Filters: []types.Filter{filterName, filterTag}})
	if err != nil {
		log.Fatal(err)
	}
	return output.Reservations
}


func DescribeEC2(outputs []types.Reservation) {

	fmt.Println("################################# EC2 Instance List #################################")
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			fmt.Printf("%s (%s): %s\n", *instance.InstanceId, instance.InstanceType, *instance.PublicDnsName)
		}
	}
}

func GetPublicDNS(outputs []types.Reservation) []string {
	dnsNames := make([]string, 0)
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			dnsNames = append(dnsNames, *instance.PublicDnsName)
		}
	}
	return dnsNames
}
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/a7420174/awscp"
)

var (
	name   string
	tagKey string
)

func init() {
	flag.StringVar(&name, "name", "", "Name of EC2 instance")
	flag.StringVar(&tagKey, "tag-key", "", "Tag key of EC2 instance")
}

func main() {
	flag.Parse()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Println("################################# EC2 Instance List #################################")
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("%s (%s): %s\n", *instance.InstanceId, instance.InstanceType, *instance.PublicDnsName)
		}
	}



}

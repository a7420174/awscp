package awscp

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

// GetReservations returns a list of reservations
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

	outputs, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{Filters: []types.Filter{filterName, filterTag}})
	if err != nil {
		log.Fatal(err)
	}
	return outputs.Reservations
}

// DescribeEC2 prints the ids of the instances and the public DNS names
func DescribeEC2(outputs []types.Reservation) {

	fmt.Println("################################# EC2 Instance List #################################")
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			fmt.Printf("%s (%s): %s\n", *instance.InstanceId, instance.InstanceType, *instance.PublicDnsName)
		}
	}
}

// GetPublicDNSName returns the public DNS names of the instances
func GetPublicDNS(outputs []types.Reservation) []string {
	dnsNames := make([]string, 0)
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			dnsNames = append(dnsNames, *instance.PublicDnsName)
		}
	}
	return dnsNames
}

// GetPlatformName returns the platform name of the instances
func GetPlatformName(outputs []types.Reservation) []string {
	platformNames := make([]string, 0)
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			platformNames = append(platformNames, *instance.PlatformDetails)
		}
	}
	return platformNames
}

// GetImageId returns the image ids of the instances	
func GetImageId(outputs []types.Reservation) []string {
	imageIds := make([]string, 0)
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			imageIds = append(imageIds, *instance.ImageId)
		}
	}
	return imageIds
}


func GetPlatformDetails(cfg aws.Config, imageIds []string) []string {
	client := ec2.NewFromConfig(cfg)

	outputs, err := client.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{ImageIds: imageIds})
	if err != nil {
		log.Fatal(err)
	}

	platforms := make([]string, 0)
	for _, image := range outputs.Images {
		platforms = append(platforms, *image.PlatformDetails)
	}
	return platforms
}
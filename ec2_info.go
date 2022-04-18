package awscp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var (
	errNoInfo             = errors.New("can't predict the platform: No platform info in the image description")
	errMultiImage         = errors.New("can't predict the platform: Multiple images used for instances found")
	errNoRunningInstances = errors.New("no running instances found")
)

// GetReservations returns a list of reservations
func GetReservations(cfg aws.Config, name string, tagKey string, running bool) []types.Reservation {
	client := ec2.NewFromConfig(cfg)

	var filterName, filterTag, filterStatus types.Filter
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

	if running {
		tag3 := "instance-state-name"
		filterStatus = types.Filter{
			Name:   &tag3,
			Values: []string{"running"},
		}
	}

	outputs, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{Filters: []types.Filter{filterName, filterTag, filterStatus}})
	if err != nil {
		log.Fatal(err)
	}
	return outputs.Reservations
}

// DescribeEC2 prints the ids and states of the instances
func DescribeEC2(outputs []types.Reservation) {
	states := make([]types.InstanceStateName, 0)
	fmt.Println("################################# EC2 Instance List #################################")
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			fmt.Printf("%s (%s): %v\n", *instance.InstanceId, instance.InstanceType, instance.State.Name)
			states = append(states, instance.State.Name)
		}
	}
	for _, state := range states {
		if state == "running" {
			return
		}
	}
	log.Fatal(errNoRunningInstances)
}

// GetInstanceIds returns the instance ids of the instances
func GetInstanceId(outputs []types.Reservation) []string {
	instacneIds := make([]string, 0)
	for _, reservation := range outputs {
		for _, instance := range reservation.Instances {
			instacneIds = append(instacneIds, *instance.InstanceId)
		}
	}
	return instacneIds
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

func GetImageDescription(cfg aws.Config, imageIds []string) []string {
	client := ec2.NewFromConfig(cfg)

	outputs, err := client.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{ImageIds: imageIds})
	if err != nil {
		log.Fatal(err)
	}

	info := make([]string, 0)
	for _, image := range outputs.Images {
		info = append(info, *image.Description)
	}
	return info
}

func PredictPlatform(info []string) string {
	if len(info) > 1 {
		log.Fatal(errMultiImage)
	}
	var platform string
	m1, _ := regexp.MatchString("(?i)"+"Amazon Linux", info[0])
	m2, _ := regexp.MatchString("(?i)"+"Ubuntu", info[0])
	m3, _ := regexp.MatchString("(?i)"+"CentOS", info[0])
	m4, _ := regexp.MatchString("(?i)"+"Red Hat", info[0])
	m5, _ := regexp.MatchString("(?i)"+"Debian", info[0])
	m6, _ := regexp.MatchString("(?i)"+"SUSE", info[0])
	switch {
	case m1:
		platform = "amazonlinux"
	case m2:
		platform = "ubuntu"
	case m3:
		platform = "centos"
	case m4:
		platform = "rhel"
	case m5:
		platform = "debian"
	case m6:
		platform = "suse"
	default:
		log.Fatal(errNoInfo)
	}
	return platform
}

func GetUsername(platform string) string {
	var username string
	switch platform {
	case "amazonlinux":
		username = "ec2-user"
	case "ubuntu":
		username = "ubuntu"
	case "centos":
		username = "ec2-user"
	case "rhel":
		username = "ec2-user"
	case "debian":
		username = "admin"
	case "suse":
		username = "ec2-user"
	default:
		username = "ec2-user"
	}
	return username
}
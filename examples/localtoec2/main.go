package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/a7420174/awscp"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	name   string
	tagKey string
	platfrom string
)

func init() {
	flag.StringVar(&name, "name", "", "Name of EC2 instance")
	flag.StringVar(&tagKey, "tag-key", "", "Tag key of EC2 instance")
	flag.StringVar(&platfrom, "platfrom", "", "OS platform of EC2 instance: amazonlinux, ubuntu, centos, rhel, debian, suse\nif empty, the platform will be predicted")
}

func main() {
	flag.Parse()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	reservations := awscp.GetReservations(cfg, name, tagKey)
	awscp.DescribeEC2(reservations)

	// dnsNames := awscp.GetPublicDNS(reservations)
	if platfrom == "" {
		imageIds := awscp.GetImageId(reservations)
		info := awscp.GetImageDescription(cfg, imageIds)
		platfrom = awscp.GetPlatform(info)
	}

	fmt.Println("Platform:", platfrom)
	
}

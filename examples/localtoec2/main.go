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

	reservations := awscp.GetReservations(cfg, name, tagKey)
	awscp.DescribeEC2(reservations)

	// dnsNames := awscp.GetPublicDNS(reservations)
	imageIds := awscp.GetImageId(reservations)
	platforms := awscp.GetPlatformDetails(cfg, imageIds)

	for _, os := range platforms {
		fmt.Println(os)
	}

}

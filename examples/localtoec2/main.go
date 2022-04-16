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

	outputs := awscp.GetReservations(cfg, name, tagKey)
	awscp.DescribeEC2(outputs)

	dnsNames := awscp.GetPublicDNS(outputs)

	for _, dnsName := range dnsNames {
		fmt.Println(dnsName)
	}
}

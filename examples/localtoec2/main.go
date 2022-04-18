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
	keyPath string
	destPath string
	filePath string
	remotePath string
	permission string
)

func init() {
	flag.StringVar(&name, "name", "", "Name of EC2 instance")
	flag.StringVar(&tagKey, "tag-key", "", "Tag key of EC2 instance")
	flag.StringVar(&platfrom, "platfrom", "", "OS platform of EC2 instance: amazonlinux, ubuntu, centos, rhel, debian, suse\nif empty, the platform will be predicted")
	flag.StringVar(&keyPath, "key-path", "", "Path of key pair")
	flag.StringVar(&destPath, "dest-path", "", "Path of destination")
	flag.StringVar(&filePath, "file-path", "", "Path of file to be copied")
	flag.StringVar(&remotePath, "remote-path", "", "Path of remote file to be saved")
	flag.StringVar(&permission, "permission", "0755", "Permission of remote file: default 0755")
}

func errhandler(dryrun bool) {
	if dryrun {
		log.Println("Dry run, Skip error handling")
		return
	}
	if keyPath == "" {
		log.Fatal("Key path is empty")
	}
	if destPath == "" {
		log.Fatal("Destination path is empty")
	}
	if filePath == "" {
		log.Fatal("File path is empty")
	}
	if remotePath == "" {
		log.Fatal("Remote path is empty")
	}
}

func main() {
	flag.Parse()
	errhandler(true)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	reservations := awscp.GetReservations(cfg, name, tagKey, false)
	
	awscp.DescribeEC2(reservations)

	reservationsRunning := awscp.GetReservations(cfg, name, tagKey, true)

	if platfrom == "" {
		imageIds := awscp.GetImageId(reservationsRunning)
		info := awscp.GetImageDescription(cfg, imageIds)
		platfrom = awscp.PredictPlatform(info)
	}
	fmt.Println("Platform:", platfrom)

	instanceIds := awscp.GetInstanceId(reservationsRunning)
	dnsNames := awscp.GetPublicDNS(reservationsRunning)
	username := awscp.GetUsername(platfrom)

	for i := range instanceIds {
		instanceId := instanceIds[i]
		dnsName := dnsNames[i]
		client := awscp.ConnectEC2(instacneId, dnsName, username, keyPath)
		fmt.Println("Connected to", client.Host)
		// awscp.CopyFile(cfg, instanceIds[i], dnsNames[i], username, keyPath, destPath, filePath, remotePath, permission)
	}
	
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

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
	if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Invalid key path")
	}
	if destPath == "" {
		log.Fatal("Destination path is empty")
	}
	if filePath == "" {
		log.Fatal("File path is empty")
	}
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Invalid file path")
	}
}

func main() {
	flag.Parse()
	errhandler(false)

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

	var wg sync.WaitGroup

	for i := range instanceIds {
		wg.Add(1)
		// client := awscp.ConnectEC2(instanceId, dnsName, username, keyPath)
		// fmt.Println("Connected to", client.Host)
		go func(i int) {
			defer wg.Done()
			awscp.CopyLocaltoEC2(instanceIds[i], dnsNames[i], username, keyPath, filePath, destPath, permission)
		}(i)
	}
	wg.Wait()
}

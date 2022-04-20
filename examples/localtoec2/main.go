package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/a7420174/awscp"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	name       string
	tagKey     string
	ids        string
	platfrom   string
	keyPath    string
	remoteDir   string
	permission string
)

func init() {
	flag.StringVar(&name, "name", "", "Name of EC2 instances")
	flag.StringVar(&tagKey, "tag-key", "", "Tag key of EC2 instances")
	flag.StringVar(&ids, "instance-ids", "", "EC2 instance IDs: e.g. i-1234567890abcdef0,i-1234567890abcdef1")
	flag.StringVar(&platfrom, "platfrom", "", "OS platform of EC2 instances: amazonlinux, ubuntu, centos, rhel, debian, suse\nif empty, the platform will be predicted")
	flag.StringVar(&keyPath, "key-path", "", "Path of key pair")
	flag.StringVar(&remoteDir, "remote-dir", "", "Path of remote directory where files are copied: default - home directory, e.g. /home/ec2-user/dir = dir")
	flag.StringVar(&permission, "permission", "0755", "Permission of remote file: default - 0755")
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
	// if remoteDir == "" {
	// 	log.Fatal("Destination path is empty")
	// }
	if flag.Arg(0) == "" {
		log.Fatal("File path is empty")
	}
	for _, filePath := range flag.Args() {
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			log.Fatal("Invalid file path")
		}
	}
}

func main() {
	// Custom usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [file1] [file2] ...\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "[file1] [file2] ...: File to be copied\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	errhandler(false)
	fmt.Print("\n")

	files := flag.Args()
	ids_slice := strings.Split(ids, ",")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	reservations := awscp.GetReservations(cfg, name, tagKey, ids_slice, false)

	awscp.DescribeEC2(reservations)

	reservationsRunning := awscp.GetReservations(cfg, name, tagKey, ids_slice, true)

	if platfrom == "" {
		imageIds := awscp.GetImageId(reservationsRunning)
		info := awscp.GetImageDescription(cfg, imageIds)
		platfrom = awscp.PredictPlatform(info)
	}
	fmt.Print("\n")
	fmt.Println("OS Platform:", platfrom)

	instanceIds := awscp.GetInstanceId(reservationsRunning)
	dnsNames := awscp.GetPublicDNS(reservationsRunning)
	username := awscp.GetUsername(platfrom)

	fmt.Print("\n")
	var wg sync.WaitGroup
	for i := range instanceIds {
		wg.Add(1)
		// client := awscp.ConnectEC2(instanceId, dnsName, username, keyPath)
		// fmt.Println("Connected to", client.Host)
		go func(i int) {
			defer wg.Done()
			for _, filePath := range files {
				awscp.CopyLocaltoEC2(instanceIds[i], dnsNames[i], username, keyPath, filePath, remoteDir, permission)
			}
		}(i)
	}
	wg.Wait()
}

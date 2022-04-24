package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	remoteDir  string
	remoteFile string
)

func init() {
	flag.StringVar(&name, "name", "", "Name of EC2 instances")
	flag.StringVar(&tagKey, "tag-key", "", "Tag key of EC2 instances")
	flag.StringVar(&ids, "instance-ids", "", "EC2 instance IDs: e.g. i-1234567890abcdef0,i-1234567890abcdef1")
	flag.StringVar(&platfrom, "platfrom", "", "OS platform of EC2 instances: amazonlinux, ubuntu, centos, rhel, debian, suse\nif empty, the platform will be predicted")
	flag.StringVar(&keyPath, "key-path", "", "Path of key pair")
	flag.StringVar(&remoteDir, "remote-dir", "", "Path of remote directory which files are copied from: relative path - home directory, e.g. /home/{username}/dir = dir")
	flag.StringVar(&remoteFile, "remote-file", "", "Path of remote file: relative path - home directory, e.g. /home/{username}/file.txt = file.txt")
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
	if len(flag.Args()) > 1 {
		log.Fatal("Too many arguments")
	}
	if _, err := os.Stat(flag.Arg(0)); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Invalid directory path")
	}
}

func main() {
	// Custom usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [local-dir]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "[local-dir]: directory path which files is copied to\nkey-path, local-dir must be specified\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	errhandler(false)
	fmt.Print("\n")

	dirPath := flag.Arg(0)
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
			if remoteFile != "" {
				client := awscp.ConnectEC2(instanceIds[i], dnsNames[i], username, keyPath)
				defer client.Close()
				awscp.CopyEC2toLocal(client, instanceIds[i], remoteFile, dirPath)
			} else {
				dirList := strings.Split(awscp.EC2RunCommand(instanceIds[i], dnsNames[i], username, keyPath, "find "+remoteDir+" -type d| grep -v '/\\.'", false), "\n")
				remoteSubDirs := dirList[1:]
				for _, remoteSubDir := range remoteSubDirs {
					remoteSubDir = strings.Replace(remoteSubDir, dirList[0]+"/", "", 1)
					os.Mkdir(filepath.Join(dirPath, instanceIds[i], remoteSubDir), 0755)
					remoteFiles := strings.Split(awscp.EC2RunCommand(instanceIds[i], dnsNames[i], username, keyPath, "find "+remoteDir+" -type f| grep -v '/\\.'", false), "\n")
					for _, filePath := range remoteFiles {
						if filePath != "" {
							filePath = strings.Replace(filePath, dirList[0]+"/", "", 1)
							client := awscp.ConnectEC2(instanceIds[i], dnsNames[i], username, keyPath)
							defer client.Close()
							awscp.CopyEC2toLocal(client, instanceIds[i], filepath.Join(dirPath, instanceIds[i], filePath))
						}
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

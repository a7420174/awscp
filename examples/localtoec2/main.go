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
	permission string
	recursive  bool
)

func init() {
	flag.StringVar(&name, "name", "", "Name of EC2 instances")
	flag.StringVar(&tagKey, "tag-key", "", "Tag key of EC2 instances")
	flag.StringVar(&ids, "instance-ids", "", "EC2 instance IDs: e.g. i-1234567890abcdef0,i-1234567890abcdef1")
	flag.StringVar(&platfrom, "platfrom", "", "OS platform of EC2 instances: amazonlinux, ubuntu, centos, rhel, debian, suse\nif empty, the platform will be predicted")
	flag.StringVar(&keyPath, "key-path", "", "Path of key pair")
	flag.StringVar(&remoteDir, "remote-dir", "", "Path of remote directory which files are copied to: default - home directory, e.g. /home/{username}/dir = dir")
	flag.StringVar(&permission, "permission", "0755", "Permission of remote file: default - 0755")
	flag.BoolVar(&recursive, "recursive", false, "Copy files recursively")
}

func errhandler(dryrun bool) {
	if dryrun {
		log.Println("Dry run, Skip error handling")
		return
	}
	if keyPath == "" {
		log.Fatalln("Key path is empty")
	}
	if _, err := os.Stat(keyPath); errors.Is(err, os.ErrNotExist) {
		log.Fatalln("Invalid key path")
	}
	if recursive {
		if flag.Arg(0) == "" {
			log.Fatalln("Directory path is empty")
		}
		if len(flag.Args()) > 1 {
			log.Fatalln("Too many arguments")
		}
		if _, err := os.Stat(flag.Arg(0)); errors.Is(err, os.ErrNotExist) {
			log.Fatalln("Directory does not exist")
		}
	} else {
		if flag.Arg(0) == "" {
			log.Fatalln("File path is empty")
		}
		for _, filePath := range flag.Args() {
			if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
				log.Fatalln("File does not exist:", filePath)
			}
		}
	}
}

func main() {
	// Custom usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [File1] [File2] ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s [flags] -recursive [Dir]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "[File1] [File2] ...: File path to be copied\n")
		fmt.Fprintf(os.Stderr, "[Dir] ...: Directory path to be copied\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	errhandler(false)
	fmt.Print("\n")

	var files []string
	var absDir string
	if recursive {
		absDir, _ = filepath.Abs(flag.Arg(0))
		files, _ = func(root string) ([]string, error) {
			var files []string
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if strings.HasPrefix(info.Name(), ".") && info.IsDir() {
					return filepath.SkipDir
				}
				if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
					files = append(files, path)
				}
				return nil
			})
			return files, err
		}(absDir)
	} else {
		files = make([]string, 0)
		for _, filePath := range flag.Args() {
			absPath, _ := filepath.Abs(filePath)
			files = append(files, absPath)
		}
	}
	// for _, filePath := range files {
	// 	fmt.Println("File:", filePath)
	// }

	ids_slice := strings.Split(ids, ",")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
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
		// client := awscp.ConnectEC2(instanceId, dnsName, username, keyPath)
		// fmt.Println("Connected to", client.Host)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for _, filePath := range files {
				client := awscp.ConnectEC2(instanceIds[i], dnsNames[i], username, keyPath)
				defer client.Close()
				if recursive {
					awscp.CopyLocaltoEC2(client, instanceIds[i], filePath, strings.Replace(filePath, absDir+"/", "", 1), permission)
				} else {
					awscp.CopyLocaltoEC2(client, instanceIds[i], filePath, filepath.Join(remoteDir, filepath.Base(filePath)), permission)
				}
			}
		}(i)
	}
	wg.Wait()
}

package awscp

import (
	"bytes"
	"context"
	"log"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// EC2RunCommand runs a command on an EC2 instance
func EC2RunCommand(instanceId, dnsName, username, keyPath, command string, verbose bool) (*bytes.Buffer, *bytes.Buffer) {
	// Connect to EC2 instance
	clientConfig, _ := auth.PrivateKey(username, keyPath, ssh.InsecureIgnoreHostKey())
	client, errDial := ssh.Dial("tcp", dnsName+":22", &clientConfig)
	if errDial != nil {
		log.Println("Error while running command ", errDial, "["+instanceId+"]")
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// Run command
	session, errSession := client.NewSession()
	if errSession != nil {
		log.Println("Error while running command ", errSession, "["+instanceId+"]")
	}

	defer session.Close()

	var b1 bytes.Buffer // import "bytes"
	var b2 bytes.Buffer // import "bytes"

	session.Stdout = &b1 // get output

	session.Stderr = &b2 // get error
	// you can also pass what gets input to the stdin, allowing you to pipe
	// content from client to server
	//      session.Stdin = bytes.NewBufferString("My input")

	// Finally, run the command
	errRun := session.Run(command)
	if errRun != nil {
		log.Println("Error while running command ", errRun, "["+instanceId+"]")
	}

	if verbose && errDial == nil && errSession == nil && errRun == nil {
		log.Println("Command executed successfully", "["+instanceId+"]")
	}
	return &b1, &b2
}

// func GetFilesRecursive(instanceId, dnsName, username, keyPath, remoteDir string) []string {
// 	var remoteFiles []string
// 	var getFiles func(dirPath string)
// 	var wg sync.WaitGroup
// 	defer wg.Wait()
// 	getFiles = func(dirPath string) {
// 		defer wg.Done()
// 		files := strings.Split(EC2RunCommand(instanceId, dnsName, username, keyPath, "ls -p "+dirPath+" | grep -v /", true), "\n")
// 		remoteFiles = append(remoteFiles, files...)
// 		directories := strings.Split(EC2RunCommand(instanceId, dnsName, username, keyPath, "ls -p "+dirPath+" | grep /", true), "\n")
// 		for _, directory := range directories {
// 			if directory != "" {
// 				wg.Add(1)
// 				go getFiles(dirPath + directory)
// 			}
// 		}
// 	}
// 	wg.Add(1)
// 	go getFiles(remoteDir)

// 	return remoteFiles
// }

// ConnectEC2 connects to an EC2 instance using SSH
func ConnectEC2(instanceId, dnsName, username, keyPath string) *scp.Client {
	clientConfig, _ := auth.PrivateKey(username, keyPath, ssh.InsecureIgnoreHostKey())

	client := scp.NewClient(dnsName+":22", &clientConfig)

	err := client.Connect()
	if err != nil {
		log.Fatalln("Couldn't establish a connection to the remote server ", "["+instanceId+"]")
	}

	return &client
}

// CopyLocalToEC2 copies a local file to an EC2 instance
func CopyLocaltoEC2(client *scp.Client, instanceId, filePath, remotePath, permission string) {
	// Open a file
	f, _ := os.Open(filePath)

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	err := client.CopyFromFile(context.TODO(), *f, remotePath, permission)

	if err != nil {
		log.Fatalln("Error while copying file ", err)
	}

	log.Println("File "+"("+filePath+")"+" copied successfully", "["+instanceId+"]")
}

// CopyEC2ToLocal copies a file from an EC2 instance to local
func CopyEC2toLocal(client *scp.Client, instanceId, filePath, remotePath string) {
	// Open a file
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		log.Fatalln("Couldn't open the output file:", err)
	}

	// Close the file after it has been copied
	defer f.Close()

	err = client.CopyFromRemote(context.TODO(), f, remotePath)

	if err != nil {
		log.Fatalln("Error while copying file:", err)
	}

	log.Println("File "+"("+remotePath+")"+" copied successfully", "["+instanceId+"]")
}

package awscp

import (
	"bytes"
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// ConnectEC2 connects to an EC2 instance using SSH
func ConnectEC2(instacneId, dnsName, username, keypath string) *scp.Client {
	clientConfig, _ := auth.PrivateKey(username, keypath, ssh.InsecureIgnoreHostKey())

	client := scp.NewClient(dnsName+":22", &clientConfig)

	err := client.Connect()
	if err != nil {
		log.Fatalln("Couldn't establish a connection to the remote server ", "["+instacneId+"]")
	}

	return &client
}

func EC2RunCommand(instacneId, dnsName, username, keypath, command string, verbose bool) string {
	// Connect to EC2 instance
	clientConfig, _ := auth.PrivateKey(username, keypath, ssh.InsecureIgnoreHostKey())
	client, err := ssh.Dial("tcp", dnsName+":22", &clientConfig)
	if err != nil {
		log.Fatalln("Error while running command ", err)
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// Run command
	session, err := client.NewSession()
	if err != nil {
		log.Fatalln("Error while running command ", err)
	}

	defer session.Close()

	var b bytes.Buffer  // import "bytes"
	session.Stdout = &b // get output
	// you can also pass what gets input to the stdin, allowing you to pipe
	// content from client to server
	//      session.Stdin = bytes.NewBufferString("My input")

	// Finally, run the command
	err = session.Run(command)
	if err != nil {
		log.Fatalln("Error while running command ", err)
	}

	if verbose {
		log.Println("Command executed successfully", "["+instacneId+"]")
	}
	return b.String()
}

// CopyLocalToEC2 copies a local file to an EC2 instance
func CopyLocaltoEC2(client *scp.Client, instacneId, filePath, remoteDir, permission string) {
	filename := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
	matched, _ := regexp.MatchString("/$", remoteDir)
	var remotePath string
	if remoteDir == "" || matched {
		remotePath = remoteDir + filename
	} else {
		remotePath = remoteDir + "/" + filename
	}

	// Open a file
	f, _ := os.Open(filePath)
	
	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFromFile(context, file, remotePath, permission)

	err := client.CopyFromFile(context.TODO(), *f, remotePath, permission)

	if err != nil {
		log.Fatalln("Error while copying file ", err)
	}

	log.Println("File "+"("+filename+")"+" copied successfully", "["+instacneId+"]")
}

// CopyEC2ToLocal copies a file from an EC2 instance to local
func CopyEC2toLocal(client *scp.Client, instacneId, filePath, localDir string) {
	filename := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
	matched, _ := regexp.MatchString("/$", localDir)
	var localPath string
	if matched {
		os.Mkdir(localDir+instacneId, 0775)
		localPath = localDir + instacneId + "/" + filename
	} else {
		os.Mkdir(localDir+"/"+instacneId, 0775)
		localPath = localDir + "/" + instacneId + "/" + filename
	}
	// Open a file
	f, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		log.Fatalln("Couldn't open the output file:", err)
	}

	// Close the file after it has been copied
	defer f.Close()

	err = client.CopyFromRemote(context.TODO(), f, filePath)

	if err != nil {
		log.Fatalln("Error while copying file:", err)
	}

	log.Println("File "+"("+filename+")"+" copied successfully", "["+instacneId+"]")
}
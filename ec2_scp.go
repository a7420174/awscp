package awscp

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func ConnectEC2(instacneId, dnsName, username, keypath string) *scp.Client {
	clientConfig, _ := auth.PrivateKey(username, keypath, ssh.InsecureIgnoreHostKey())

	client := scp.NewClient(dnsName+":22", &clientConfig)

	err := client.Connect()
	if err != nil {
		log.Println("Couldn't establish a connection to the remote server ", "["+instacneId+"]")
	}

	return &client
}

func CopyLocaltoEC2(instacneId, dnsName, username, keypath, filepath, remoteDir, permission string) {
	// Connect to EC2 instance
	client := ConnectEC2(instacneId, dnsName, username, keypath)

	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	matched, _ := regexp.MatchString("/$", remoteDir)
	var remotePath string
	if remoteDir == "" || matched {
		remotePath = remoteDir + filename
	} else {
		remotePath = remoteDir + "/" + filename
	}

	// Open a file
	f, _ := os.Open(filepath)

	// Close client connection after the file has been copied
	defer client.Close()

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

func CopyEC2toLocal(instacneId, dnsName, username, keypath, filepath, localDir string) {
	// Connect to EC2 instance
	client := ConnectEC2(instacneId, dnsName, username, keypath)

	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
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
		log.Fatal("Couldn't open the output file")
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	err = client.CopyFromRemote(context.TODO(), f, filepath)

	if err != nil {
		log.Fatalln("Error while copying file ", err)
	}

	log.Println("File "+"("+filename+")"+" copied successfully", "["+instacneId+"]")
}

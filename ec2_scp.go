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

func CopyLocaltoEC2(instacneId, dnsName, username, keypath, filepath, destpath, permission string) {
	// Connect to EC2 instance
	client := ConnectEC2(instacneId, dnsName, username, keypath)

	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	matched, _ := regexp.MatchString("/$", destpath)
	if destpath == "" || matched {
		destpath = destpath + filename
	}

	// Open a file
	f, _ := os.Open(filepath)

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFromFile(context, file, destpath, permission)

	err := client.CopyFromFile(context.TODO(), *f, destpath, permission)

	if err != nil {
		log.Println("Error while copying file ", err)
	}

	log.Println("File ", "("+filename+")", " copied successfully", "["+instacneId+"]")
}

func CopyEC2toLocal(instacneId, dnsName, username, keypath, filepath, destpath, permission string) {
	// Connect to EC2 instance
	client := ConnectEC2(instacneId, dnsName, username, keypath)

	filename := strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	matched, _ := regexp.MatchString("/$", destpath)
	if destpath == "" || matched {
		destpath = destpath + filename
	}

	// Open a file
	f, _ := os.Open(destpath)

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFromFile(context, file, destpath, permission)

	err := client.CopyFromRemote(context.TODO(), f, filepath)

	if err != nil {
		log.Println("Error while copying file ", err)
	}

	log.Println("File ", "("+filename+")", " copied successfully", "["+instacneId+"]")
}

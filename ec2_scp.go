package awscp

import (
	"context"
	"log"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func ConnectEC2(instacneId, dnsName, username, keypath string) *scp.Client {
	clientConfig, _ := auth.PrivateKey(username, keypath, ssh.InsecureIgnoreHostKey())

	client := scp.NewClient(dnsName+":22", &clientConfig)

	err := client.Connect()
	if err != nil {
		log.Println("Couldn't establish a connection to the remote server ", "("+instacneId+")")
	}

	return &client
}

func CopyLocaltoEC2(instacneId, dnsName, username, keypath, filepath, destpath, permission string) {
	// Connect to EC2 instance
	client := ConnectEC2(instacneId, dnsName, username, keypath)

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

	log.Println("File copied successfully", "("+instacneId+")")
}

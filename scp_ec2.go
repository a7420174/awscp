package awscp

import (
	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func ConnectEC2(host string, port int, user string, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			auth.Password(password),
		},
	}
	client, err := scp.DialWithConfig(host+":"+string(port), config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CopyFromLocal(dnsName string, username string, keypath string) {
	clientConfig, _ := auth.PrivateKey(username, , ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	
	client := scp.NewClient(+"22", &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	// Open a file
	f, _ := os.Open("/path/to/local/file")

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFromFile(context, file, remotePath, permission)

        // the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
	err = client.CopyFromFile(context.Background(), *f, "/home/server/test.txt", "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
}


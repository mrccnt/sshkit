# SSHKit

[![license](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![goreportcard](https://goreportcard.com/badge/github.com/mrccnt/sshkit)](https://goreportcard.com/report/github.com/mrccnt/sshkit)

This package encapsulates my most common SSH actions needed in daily business. It simplifies creating/handling of SSH
connections, simple SSH tunnels or SFTP filetransfers. Additionally it adds `AgentAuth` as `ssh.AuthMethod`. This way
you can use your local SSH agent for authentication and do not have to handle passwords.

## Examples

Keep in mind that these are examples. Handle your errors! Do not panic! ;)

### SFTP

Transfer a file via SFTP from remote server to locahost:

```go
package main

import (
	"github.com/mrccnt/sshkit"
	"golang.org/x/crypto/ssh"
	"fmt"
	"os"
)

func main() {

	sshClient, err := sshkit.SSHClient(
		sshkit.SSHConfig(
			"username",
			[]ssh.AuthMethod{sshkit.AgentAuth()},
		),
		"domain.tld:22",
	)
	if err != nil {
		panic(err.Error())
	}
	defer sshClient.Close()

	sftpClient, err := sshkit.SFTPClient(sshClient)
	if err != nil {
		panic(err.Error())
	}
	defer sftpClient.Close()

	bytes, err := sshkit.Pull(sftpClient, "remote/src/file.md", "local/dst/file.md")
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(bytes, "bytes written on local filesystem")

    bytes, err = sshkit.Push(sftpClient, "local/src/file.md", "remote/dst/file.md")
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(bytes, "bytes written on remote storage")
    
}
```

### Tunnel

Create a tunnel to a remote host to make remote mysql server available on localhost port 13306:

```go
package main

import (
    "github.com/mrccnt/sshkit"	
    "golang.org/x/crypto/ssh"	
)

func main(){
	// Create the complete tunnel configuration
	// You need a ssh.ClientConfig and someEndpoints
	tunnel := &sshkit.SSHTunnel{
		// Create a SSH ClientConfiguration with your requirements.
		Config: &ssh.ClientConfig{
			User: "username",
			Auth: []ssh.AuthMethod{
				sshkit.AgentAuth(),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		// The tunneled port from your remote server will be available on your localhost on port 13306
		Local: &sshkit.Endpoint{
			Host: "localhost",
			Port: 13306,
		},
		// Your remote SSH server
		Server: &sshkit.Endpoint{
			Host: "domain.tld",
			Port: 22,
		},
		// On your remote host we will use localhosts MySql server on default port 3306
		Remote: &sshkit.Endpoint{
			Host: "localhost",
			Port: 3306,
		},
	}
	// Initialize/Start the tunnel.
	// Keep in mind: This is a blocking action...
	err := tunnel.Start()
	if err != nil {
		panic(err.Error())
	}
}
```
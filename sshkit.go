// Copyright 2017 Marco Conti
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sshkit

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
)

type (
	// This is the config needed to create a tunnel
	SSHTunnel struct {
		Local  *Endpoint
		Server *Endpoint
		Remote *Endpoint
		Config *ssh.ClientConfig
	}
	// Endpoint defines a host and port
	Endpoint struct {
		Host string
		Port int
	}
)

// SSHConfig returns a ssh.ClientConfig for given parameters
func SSHConfig(username string, auths []ssh.AuthMethod) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// SSHClient returns a SSH Client for given config and address
func SSHClient(sshCfg *ssh.ClientConfig, addr string) (*ssh.Client, error) {
	c, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// LocalAgent returns a connection to the local authentication agent via auth sock env variable
func LocalAgent() (net.Conn, error) {
	a, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}
	return a, nil
}

// AgentAuth returns a ssh.AuthMethod to use your local authentication agent
func AgentAuth() ssh.AuthMethod {
	a, err := LocalAgent()
	if err != nil {
		panic(err.Error())
	}
	return ssh.PublicKeysCallback(agent.NewClient(a).Signers)
}

// SFTPClient returns a client to handle any SFTP actions for given SSH Client
func SFTPClient(sshClient *ssh.Client) (*sftp.Client, error) {
	c, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, nil
	}
	return c, nil
}

// String returns the preferred ssh package address format
func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// Start initializes the tunnel and starts a forwarding (blocking)
func (tunnel *SSHTunnel) Start() error {
	l, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer l.Close()
	for {
		local, err := l.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(local)
	}
}

// forward handles any incoming/outgoing reader/writer for this tunnel
func (tunnel *SSHTunnel) forward(local net.Conn) {

	server, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remote, err := server.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(local, remote)
	go copyConn(remote, local)
}

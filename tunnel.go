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
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os/exec"
)

// Endpoint defines a host and port
type Endpoint struct {
	Host string
	Port int
}

// String returns the preferred ssh package address format
func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// SSHTunnel is the config needed to create a tunnel
type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
	socket string
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

// StartCli creates tunnel to/from remote host via shell exec using ssh commands
func (tunnel *SSHTunnel) StartCli() error {

	ident, err := uuid.NewV4()
	if err != nil {
		panic(err.Error())
	}

	tunnel.socket = ident.String()

	var cmd *exec.Cmd

	cmd = exec.Command(
		"ssh",
		"-p",
		string(tunnel.Server.Port),
		"-M",
		"-S",
		tunnel.socket,
		"-fnNT",
		"-L",
		tunnel.route(),
		tunnel.userhost(),
	)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err

	}

	return nil
}

// StopCli closes tunnel to/from remote host via shell exec using ssh commands
func StopCli(tunnel *SSHTunnel) error {

	if tunnel.socket == "" {
		return errors.New("No socket identifier available")
	}

	var cmd *exec.Cmd

	cmd = exec.Command(
		"ssh",
		"-p",
		string(tunnel.Server.Port),
		"-M",
		"-S",
		tunnel.socket,
		"-O",
		"exit",
		tunnel.userhost(),
	)

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
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

// userhost returns preformated string of username and host/ip used in CLI actions
func (tunnel *SSHTunnel) userhost() string {
	return tunnel.Config.User + "@" + tunnel.Server.Host
}

// route returns preformated string of local port, remote hostname and remote port used in CLI actions
func (tunnel *SSHTunnel) route() string {
	return fmt.Sprintf(
		"%s:%s:%s",
		string(tunnel.Local.Port),
		tunnel.Remote.Host,
		string(tunnel.Remote.Port),
	)
}

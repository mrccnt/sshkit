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
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
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

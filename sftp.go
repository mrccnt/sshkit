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
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"reflect"
)

// SFTPClient returns a client to handle any SFTP actions
func SFTPClient(sshClient *ssh.Client) (*sftp.Client, error) {
	c, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, nil
	}
	return c, nil
}

// Push uploads a file from local fs to remote server
func Push(client *sftp.Client, local string, remote string) (int64, error) {

	srcFile, err := os.Open(local)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(remote)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// Pull downloads a file from remote server to local fs
func Pull(client *sftp.Client, remote string, local string) (int64, error) {

	srcFile, err := client.Open(remote)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(local)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return srcFile.WriteTo(dstFile)
}

// IsReadable checks if a given file on the remote server is readable
func IsReadable(client *sftp.Client, path string) bool {
	f, err := client.Open(path)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

// Exists checks via glob() if a given file exists on the remote server
func Exists(client *sftp.Client, path string) (bool, error) {
	results, err := client.Glob(path)
	if err != nil {
		return false, err
	}
	if len(results) == 1 && results[0] == path {
		return true, nil
	}
	return false, nil
}

// HasOsAttrib checks if given os.FileInfo matches given permissions:
//
// Use OS_ constants to define the permissions flag:
//
// if sshkit.HasOsAttrib(info, sshkit.OsUserR) {
// 		// This file is readable for user
// }
//
// if !sshkit.HasOsAttrib(info, sshkit.OsGroupW) {
// 		// This file is NOT writable for group members
// }
//
func HasOsAttrib(info os.FileInfo, permission uint32) bool {
	fs := reflect.ValueOf(info.Sys()).Elem().Interface().(sftp.FileStat)
	return fs.Mode&(permission) != 0
}

// GetUID returns the UID of given os.FileInfo
func GetUID(info os.FileInfo) uint32 {
	fs := reflect.ValueOf(info.Sys()).Elem().Interface().(sftp.FileStat)
	return fs.UID
}

// GetGID returns the GID of given os.FileInfo
func GetGID(info os.FileInfo) uint32 {
	fs := reflect.ValueOf(info.Sys()).Elem().Interface().(sftp.FileStat)
	return fs.GID
}

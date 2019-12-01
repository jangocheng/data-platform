package ftp

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func connect(user, passWd string, host string, port int) (*sftp.Client, error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(passWd))
	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "fail to connect to remote host")
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, errors.Wrap(err, "fail to connect to remote host")
	}
	return sftpClient, nil
}

package ftp

import (
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
)

func SendFile(user, passWd string, host string, port int, content, remotePath string) error {
	sftpClient, err := connect(user, passWd, host, port)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	var remoteFile *sftp.File
	remoteFile, err = sftpClient.Create(remotePath)
	if err != nil {
		return errors.Wrap(err, "fail to send content to remote file")
	}
	defer remoteFile.Close()
	_, err = remoteFile.Write([]byte(content))
	if err != nil {
		return errors.Wrap(err, "fail to send content to remote file")
	}
	return nil
}

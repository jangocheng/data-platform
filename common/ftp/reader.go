
package ftp

import (
	"bufio"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"io"
	"os"
)

func ReadFile(user, passWd string, host string, port int, fileName string) (string, error) {
	sftpClient, err := connect(user, passWd, host, port)
	if err != nil {
		return "", err
	}
	defer sftpClient.Close()
	remoteFileStat, err := sftpClient.Stat(fileName)
	if err != nil {
		return "", errors.Wrap(err, "fail to read file form remote")
	}
	remoteFile, err := sftpClient.Open(fileName)
	if err != nil {
		return "", errors.Wrap(err, "fail to read file from remote")
	}
	defer remoteFile.Close()
	buffer := make([]byte, remoteFileStat.Size())
	_, err = remoteFile.Read(buffer)
	if err != nil {
		return "", errors.Wrap(err, "fail to read file from remote")
	}
	return string(buffer), err
}


func ReadLines(user, passWd string, host string, port int, fileName string) ([]string, error) {
	sftpClient, err := connect(user, passWd, host, port)
	if err != nil {
		return nil, err
	}
	var (
		file 	*sftp.File
		part 	[]byte
		prefix 	bool
		lines 	[]string
	)
	defer sftpClient.Close()
	_, err = sftpClient.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, errors.New("fail to read file lines from remote")
	}

	file, err = sftpClient.Open(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read file lines from remote")
	}

	reader := bufio.NewReader(file)
	defer file.Close()

	for {
		part, prefix, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, errors.Wrap(err, "fail to read file lines from remote")
			}
		}
		if !prefix {
			lines = append(lines, string(part))
		}
	}
	return lines, nil
}


func ReadDir(user, passWd string, host string, port int, folder string) ([]string, error) {
	sftpClient, err := connect(user, passWd, host, port)
	if err != nil {
		return nil, err
	}
	defer sftpClient.Close()
	files, err := sftpClient.ReadDir(folder)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read dir")
	}
	var dirs []string
	for _, file := range files {
		dirs = append(dirs, file.Name())
	}
	return dirs, nil
}






package utils

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
)

func OverwriteFile(content string, fileName string) error {
	var dstFile *os.File
	dstFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}


func AppendToFile(content string, fileName string) error {
	var dstFile *os.File
	dstFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}


func AppendLinesToFile(lines []string, fileName string) error {
	var dstFile *os.File
	dstFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dstFile.Close()
	for _, content := range lines {
		_, err = dstFile.WriteString(content + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}


func GetFileContent(fileName string) (string, error) {
	fileStat, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return "", errors.New("can not find file name :" + " " + fileName)
	}
	fileDst, err := os.Open(fileName)
	defer fileDst.Close()
	if err != nil {
		return "", err
	}
	buffer := make([]byte, fileStat.Size())
	_, err = fileDst.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func GetFileContentLines(fileName string)([]string, error){
	var (
		file 	*os.File
		part 	[]byte
		prefix 	bool
		lines 	[]string
	)

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return []string{}, errors.New("can not find file name :" + " " + fileName)
	}

	file, err = os.Open(fileName)
	if err != nil {
		return []string{}, err
	}

	reader := bufio.NewReader(file)
	defer file.Close()

	for {
		part, prefix, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return []string{}, err
			}
		}
		if !prefix {
			lines = append(lines, string(part))
		}
	}
	return lines, nil
}


func Mkdir(path string) error {
	err := os.MkdirAll(path, 0766)
	if err != nil {
		return err
	}
	return nil
}

func ReadDir(dir string) (result []string) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() {
			result = append(result, f.Name())
		}
	}
	return result
}

func Remove(path string) error {
	err := os.RemoveAll(path)
	return err
}

package utils

import (
	"encoding/json"
	io "io/ioutil"
)

type JsonParser struct{

}

func NewJsonParser () *JsonParser {
	return &JsonParser{}

}

func (t *JsonParser) load (filename string, cnf interface{}) error {

	data, err := io.ReadFile(filename)
	if err != nil{
		return err
	}

	err = json.Unmarshal(data, cnf)
	if err != nil{
		return err
	}
	return nil
}

func LoadJsonConf(fileName string, cnf interface{}) error {
	parser := NewJsonParser()
	err := parser.load(fileName, cnf)
	if err != nil {
		return err
	}
	return nil
}

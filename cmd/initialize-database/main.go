package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/dryairship/online-election-manager/config"
)

type InitPost struct {
	Id         string
	Name       string
	HasNota    bool
	Candidates []string
}

type InitData struct {
	Posts []InitPost
}

func main() {
	fileData, err := ioutil.ReadFile(config.ElectionDataFilePath)
	if err != nil {
		log.Fatal("[ERROR] Election Data file not found.")
	}
	var data InitData
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("[ERROR] Election Data is not in the correct format.")
	}

	jsonData := AddToDb(data)
	WriteToFile(jsonData)
	UpdateVoters(data)
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/dryairship/online-election-manager/config"
)

func WriteToFile(posts []jsonPost) {
	jsonStr, err := json.Marshal(posts)
	if err != nil {
		log.Fatal("Cannot convert posts data to json. Error:", err.Error())
	}
	err = ioutil.WriteFile(config.CandidatesOutputPath+"/candidates", jsonStr, 0644)
	if err != nil {
		log.Fatal("Cannot save json data to file. Error:", err.Error())
	}
}

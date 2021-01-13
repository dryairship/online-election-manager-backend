package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
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
	log.Println("Reading voters lists...")
	allVoters := readVotersLists()
	log.Println("Adding voters to db...")
	createVoters(allVoters)
}

func readVotersLists() []string {
	fileData, err := ioutil.ReadFile(config.ElectionDataFilePath)
	if err != nil {
		log.Fatal("[ERROR] Election Data file not found.")
	}
	var data InitData
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("[ERROR] Election Data is not in the correct format.")
	}

	var allVoters []string
	for _, post := range data.Posts {
		voters := getVotersForPost(post.Id)
		allVoters = append(allVoters, voters...)
	}

	return removeDuplicateVoters(allVoters)
}

func removeDuplicateVoters(allVoters []string) []string {
	alreadyCounted := make(map[string]bool)
	voters := []string{}
	for _, voter := range allVoters {
		if counted := alreadyCounted[voter]; !counted {
			alreadyCounted[voter] = true
			voters = append(voters, voter)
		}
	}
	return voters
}

func getVotersForPost(id string) []string {
	file, err := os.Open(config.VotersListPath + "/" + id + ".txt")
	if err != nil {
		log.Fatalf("Error while opening voters list for post <%s>: %s", id, err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var allVoters []string

	for scanner.Scan() {
		allVoters = append(allVoters, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while scanning voters list for post <%s>: %s", id, err.Error())
	}
	log.Println("Successfully read", id, ".txt")
	return allVoters
}

func createVoters(voterRolls []string) {
	electionDb, err := db.ConnectToDatabase()
	if err != nil {
		log.Fatal("[ERROR] Could not establish database connection.")
	}

	for _, roll := range voterRolls {
		skeleton, err := electionDb.FindStudentSkeleton(roll)
		if err != nil {
			log.Println(roll, "not found.")
			continue
		}
		voter := skeleton.CreateVoter("")
		err = electionDb.AddNewVoter(&voter)
		if err != nil {
			log.Println("Could not add", roll)
		}
	}
	log.Println("Added voters.")
}

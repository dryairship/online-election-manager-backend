package main

import (
	"bufio"
	"log"
	"os"

	"github.com/dryairship/online-election-manager/config"
)

func UpdateVoters(data InitData) {
	err := electionDb.ClearPostsForAllVoters()
	if err != nil {
		log.Fatalf("Error while clearing posts for all voters: %s", err.Error())
	}

	m := make(map[string][]string)

	for _, post := range data.Posts {
		file, err := os.Open(config.VotersListPath + "/" + post.Id + ".txt")
		if err != nil {
			log.Fatalf("Error while opening voters list for post <%s>: %s", post.Id, err.Error())
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var roll string
		for scanner.Scan() {
			roll = scanner.Text()
			m[roll] = append(m[roll], post.Id)
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error while scanning voters list for post <%s>: %s", post.Id, err.Error())
		}
	}

	fail := false
	for roll := range m {
		_, err = electionDb.FindVoter(roll)
		if err != nil {
			fail = true
			log.Printf("Voter <%s> not found in the database: %s\n", roll, err.Error())
		}
	}

	if fail {
		log.Fatal("One or more voters were not found in the database.")
	}

	for roll, posts := range m {
		err = electionDb.SetPostsForVoter(roll, posts)
		if err != nil {
			log.Fatalf("Error while setting posts for voter <%s>: %s", roll, err.Error())
		}
	}
}

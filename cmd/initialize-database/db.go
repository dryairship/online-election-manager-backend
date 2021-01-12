package main

import (
	"fmt"
	"log"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
)

var electionDb db.ElectionDatabase

type jsonCandidate struct {
	Roll string `json:"roll"`
	Name string `json:"name"`
}

type jsonPost struct {
	Id         string          `json:"postId"`
	Name       string          `json:"postName"`
	Candidates []jsonCandidate `json:"candidates"`
}

func AddToDb(data InitData) []jsonPost {
	var err error
	electionDb, err = db.ConnectToDatabase()
	if err != nil {
		log.Fatal("[ERROR] Could not establish database connection.")
	}
	err = electionDb.ResetDatabase()
	if err != nil {
		log.Fatal("[ERROR] Could not reset database.")
	}

	result := make([]jsonPost, len(data.Posts))

	for i, post := range data.Posts {
		var candidatesUsernames []string
		jsonCands := make([]jsonCandidate, len(post.Candidates))

		for j, candidateRoll := range post.Candidates {
			candidate := CreateCandidate(&electionDb, post.Id, candidateRoll)
			candidatesUsernames = append(candidatesUsernames, candidate.Username)
			err = electionDb.AddNewCandidate(&candidate)
			if err != nil {
				log.Fatalf("[ERROR] Cannot add candidate %s\n", candidateRoll)
			}
			jsonCands[j] = jsonCandidate{
				Roll: candidate.Roll,
				Name: candidate.Name,
			}
		}

		fullPost := models.Post{
			PostID:     post.Id,
			PostName:   post.Name,
			Candidates: candidatesUsernames,
			Resolved:   false,
		}

		result[i] = jsonPost{
			Id:         post.Id,
			Name:       post.Name,
			Candidates: jsonCands,
		}

		// Insert the newly created post into the database.
		err = electionDb.AddNewPost(&fullPost)

		if err != nil {
			log.Fatalf("[ERROR] Cannot insert post %s\n", post.Name)
		} else {
			log.Printf("Succesfully added the post of %s.\n", post.Name)
		}
	}
	return result
}

// Function to create a candidate from the data in the file.
func CreateCandidate(dB *db.ElectionDatabase, pID, roll string) models.Candidate {
	skeleton, err := dB.FindStudentSkeleton(roll)
	if err != nil {
		log.Fatalf("[ERROR] Cannot find candidate %s\n", roll)
	}

	candidate := models.Candidate{
		Roll:       roll,
		Name:       skeleton.Name,
		Email:      skeleton.Email + config.MailSuffix,
		Username:   fmt.Sprintf("P%sC%s", pID, roll),
		Password:   "",
		AuthCode:   "",
		PostID:     pID,
		PublicKey:  "",
		PrivateKey: "",
		KeyState:   config.KeysNotGenerated,
	}

	return candidate
}

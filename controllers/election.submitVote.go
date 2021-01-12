package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/models"
	"github.com/dryairship/online-election-manager/utils"
)

// API handler to store the submitted vote.
func SubmitVote(c *gin.Context) {
	roll := c.GetString("ID")
	voter, err := ElectionDb.FindVoter(roll)
	if err != nil {
		log.Println("[ERROR] Voter not found in db while submitting vote: ", roll, err.Error())
		c.String(http.StatusNotFound, "Voter not found.")
		return
	}

	if voter.Voted {
		log.Println("[WARN] Attempt to re-vote: ", voter)
		c.String(http.StatusForbidden, "This voter has already voted.")
		return
	}

	var receivedVotes []models.ReceivedVote
	err = c.BindJSON(&receivedVotes)
	if err != nil {
		log.Println("[ERROR] Vote JSON could not bind to struct: ", voter, err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	for _, receivedVote := range receivedVotes {
		isAllowed := false
		for _, pId := range voter.Posts {
			if pId == receivedVote.PostID {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			log.Printf("[WARN] Voter <%s> wrongly tried to vote on post <%s>\n", voter.Roll, receivedVote.PostID)
			c.String(http.StatusForbidden, "This voter is not allowed to vote for post "+receivedVote.PostID)
			return
		}
	}

	ballotID := make([]models.BallotID, len(receivedVotes))
	for i, receivedVote := range receivedVotes {
		ballotID[i] = receivedVote.GetBallotID()
		vote := receivedVote.GetVote()

		// Insert the vote after a random time delay.
		go func(vote models.Vote, voter models.Voter) {
			time.Sleep(utils.GetRandomTimeDelay())
			err := ElectionDb.InsertVote(&vote)
			if err != nil {
				log.Println("[ERROR] Vote could not be inserted in the database: ", voter, err.Error())
			}
		}(vote, voter)
	}

	newVoter := models.Voter{
		Roll:     voter.Roll,
		Name:     voter.Name,
		Email:    voter.Email,
		Password: voter.Password,
		AuthCode: voter.AuthCode,
		BallotID: ballotID,
		Voted:    true,
	}

	err = ElectionDb.UpdateVoter(roll, &newVoter)
	if err != nil {
		log.Println("[ERROR] Voter could not be updated after vote insertion: ", voter, newVoter, err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	c.JSON(http.StatusOK, "Votes successfully submitted.")
}

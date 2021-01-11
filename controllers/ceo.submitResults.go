package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
	"github.com/dryairship/online-election-manager/utils"
)

type candidateResult struct {
	Roll        string `json:"roll"`
	Name        string `json:"name"`
	Preference1 int32  `json:"preference1"`
	Preference2 int32  `json:"preference2"`
	Preference3 int32  `json:"preference3"`
}

type postResult struct {
	ID         string            `json:"postId"`
	Name       string            `json:"postName"`
	Candidates []candidateResult `json:"candidates"`
	BallotIDs  []string          `json:"ballotIds"`
}

type resultData struct {
	Posts []postResult `json:"posts"`
}

func SubmitResults(c *gin.Context) {
	var results resultData
	err := c.BindJSON(&results)
	if err != nil {
		log.Println("[ERROR] Results JSON did not bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	for _, post := range results.Posts {
		postId := post.ID
		result := models.Result{
			ID:   postId,
			Name: post.Name,
		}
		var candidateResults []models.CandidateResult
		for _, candidate := range post.Candidates {
			candidateResult := models.CandidateResult{
				Name:        candidate.Name,
				Roll:        candidate.Roll,
				Preference1: candidate.Preference1,
				Preference2: candidate.Preference2,
				Preference3: candidate.Preference3,
				Status:      "none",
			}
			candidateResults = append(candidateResults, candidateResult)
		}
		result.Candidates = candidateResults

		err = ElectionDb.InsertResult(&result)
		if err != nil {
			log.Println("[ERROR] Database error while inserting Result: ", result, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}

		err = utils.ExportBallotIdsToFile(post.BallotIDs, postId)
		if err != nil {
			log.Println("[ERROR] Util error while exporting ballot IDs to file: ", err.Error())
			c.String(http.StatusInternalServerError, "Cannot export BallotIds.")
			return
		}
	}

	Posts, err = ElectionDb.GetPosts()
	if err != nil {
		log.Println("[ERROR] Database error while getting remaining posts: ", err.Error())
		c.String(http.StatusBadRequest, "Database error.")
		return
	}

	config.ElectionState = config.ResultsCalculated
	log.Println("[INFO] Results stored in the database")
	c.String(http.StatusAccepted, "Results accepted.")
}

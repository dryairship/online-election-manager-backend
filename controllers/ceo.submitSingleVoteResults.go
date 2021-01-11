package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
	"github.com/dryairship/online-election-manager/utils"
)

type singleVoteCandidateResult struct {
	Roll      string   `json:"roll"`
	Name      string   `json:"name"`
	Count     int32    `json:"count"`
	BallotIDs []string `json:"ballotIds"`
	Status    string   `json:"status"`
}

type singleVotePostResult struct {
	ID         string                      `json:"postId"`
	Name       string                      `json:"postName"`
	Resolved   bool                        `json:"resolved"`
	Candidates []singleVoteCandidateResult `json:"candidates"`
}

type singleVoteResultData struct {
	Posts []singleVotePostResult `json:"posts"`
}

func SubmitSingleVoteResults(c *gin.Context) {
	var results singleVoteResultData
	err := c.BindJSON(&results)
	if err != nil {
		log.Println("[ERROR] SingleVoteResults JSON did not bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	for _, post := range results.Posts {
		var ballotIds []models.UsedSingleVoteBallotID
		postId := post.ID
		if post.Resolved {
			if err := ElectionDb.MarkPostResolved(postId); err != nil {
				log.Println("[ERROR] Database error while marking post resolved: ", postId, err.Error())
				c.String(http.StatusBadRequest, "Database error.")
				return
			}
		}
		singleVoteResult := models.SingleVoteResult{
			ID:   postId,
			Name: post.Name,
		}
		var candidateResults []models.SingleVoteCandidateResult
		for _, candidate := range post.Candidates {
			tmpBallotId := models.UsedSingleVoteBallotID{
				Name: candidate.Name,
				Roll: candidate.Roll,
			}
			for _, ballotString := range candidate.BallotIDs {
				tmpBallotId.BallotString = ballotString
				ballotIds = append(ballotIds, tmpBallotId)
			}
			candidateResult := models.SingleVoteCandidateResult{
				Name:   candidate.Name,
				Roll:   candidate.Roll,
				Count:  candidate.Count,
				Status: candidate.Status,
			}
			candidateResults = append(candidateResults, candidateResult)

			if candidate.Status == "eliminated" {
				err = ElectionDb.EliminateCandidate(postId, candidate.Roll)
				if err != nil {
					log.Println("[ERROR] Database error while eliminating candidate: ", postId, candidate, err.Error())
					c.String(http.StatusInternalServerError, "Database error.")
					return
				}
			}
		}
		singleVoteResult.Candidates = candidateResults

		err = ElectionDb.InsertSingleVoteResult(&singleVoteResult)
		if err != nil {
			log.Println("[ERROR] Database error while inserting SingleVoteResult: ", singleVoteResult, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}

		err = utils.ExportSingleVoteBallotIdsToFile(ballotIds, postId)
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

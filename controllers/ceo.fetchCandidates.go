package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct to return candidates' data to the CEO.
type fetchCandidates_candidateData struct {
	Name       string `json:"name"`
	Roll       string `json:"roll"`
	PostID     string `json:"postId"`
	PrivateKey string `json:"privateKey"`
}

func isCandidateStillInRace(postId, username string) bool {
	for _, post := range Posts {
		if post.PostID == postId && !post.Resolved {
			for _, candidate := range post.Candidates {
				if candidate == username {
					return true
				}
			}
		}
	}
	return false
}

// API handler to fetch all candidates from the database.
func FetchCandidates(c *gin.Context) {
	candidates, err := ElectionDb.GetAllCandidates()
	if err != nil {
		log.Println("[ERROR] Database error while fetching candidates: ", err.Error())
		c.String(http.StatusInternalServerError, "Error while fetching candidates.")
	}

	var data []fetchCandidates_candidateData
	for _, candidate := range candidates {
		if isCandidateStillInRace(candidate.PostID, candidate.Username) {
			data = append(data, fetchCandidates_candidateData{
				Name:       candidate.Name,
				Roll:       candidate.Roll,
				PostID:     candidate.PostID,
				PrivateKey: candidate.PrivateKey,
			})
		}
	}
	c.JSON(http.StatusOK, data)
}

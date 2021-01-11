package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API handler to fetch all votes from the database.
func FetchVotes(c *gin.Context) {
	votes, err := ElectionDb.GetVotes()
	if err != nil {
		log.Println("[ERROR] Database error while fetching votes: ", err.Error())
		c.String(http.StatusInternalServerError, "Error while fetching votes.")
		return
	}
	c.JSON(http.StatusOK, &votes)
}

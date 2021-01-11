package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
)

func PrepareForNextRound(c *gin.Context) {
	err := ElectionDb.MarkAllVotersUnvoted()
	if err != nil {
		log.Println("[ERROR] Database error while marking all voters unvoted: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	err = ElectionDb.DeleteAllVotes()
	if err != nil {
		log.Println("[ERROR] Database error while deleting all votes: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	config.ElectionState = config.VotingNotYetStarted

	log.Println("[INFO] Ready for next round")
	c.String(http.StatusOK, "Ready for next round.")
}

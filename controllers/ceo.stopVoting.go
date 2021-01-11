package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
)

// API handler to stop voting process.
func StopVoting(c *gin.Context) {
	config.ElectionState = config.VotingStopped
	log.Println("[INFO] Voting stopped")
	c.String(http.StatusOK, "Voting Stopped")
}

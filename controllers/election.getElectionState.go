package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
)

// API handler to fetch the current election state.
func GetElectionState(c *gin.Context) {
	c.String(http.StatusOK, config.ElectionState)
}

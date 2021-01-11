package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API handler to get information about a candidate.
func GetCandidateInfo(c *gin.Context) {
	username := c.Param("username")

	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusNotFound, "This student is not contesting elections.")
		return
	}

	simplifiedCandidate := candidate.Simplify()

	c.JSON(http.StatusOK, simplifiedCandidate)
}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSingleVoteResults(c *gin.Context) {
	result, err := ElectionDb.FindAllSingleVoteResults()
	if err != nil {
		log.Println("[ERROR] Could not get single vote results from db: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	c.JSON(http.StatusOK, result)
}

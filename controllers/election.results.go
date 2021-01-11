package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetResults(c *gin.Context) {
	result, err := ElectionDb.FindAllResults()
	if err != nil {
		log.Println("[ERROR] Could not get results from db: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
		return
	}

	c.JSON(http.StatusOK, result)
}

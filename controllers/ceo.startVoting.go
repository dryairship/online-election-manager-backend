package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
)

// Struct to accept the keys from the client.
type newCEOData struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// API handler to start voting process by accepting CEO's public and private keys.
func StartVoting(c *gin.Context) {
	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		log.Println("[ERROR] No CEO in the database error on starting voting: ", err.Error())
		c.String(http.StatusInternalServerError, "CEO not yet assigned.")
		return
	}
	newData := newCEOData{}
	err = c.BindJSON(&newData)
	if err != nil {
		log.Println("[ERROR] CEO data JSON did not bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	ceo.PublicKey = newData.PublicKey
	ceo.PrivateKey = newData.PrivateKey
	err = ElectionDb.UpdateCEO(&ceo)
	if err != nil {
		log.Println("[ERROR] Database error while starting voting: ", err.Error())
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	config.PublicKeyOfCEO = ceo.PublicKey
	config.ElectionState = config.AcceptingVotes

	log.Println("[INFO] Voting started")
	c.String(http.StatusOK, "Voting Started.")
}

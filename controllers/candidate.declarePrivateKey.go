package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/utils"
)

// Struct to accept new keys from the client.
type declarePrivateKey_Keys struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

// API handler to accept the unencrypted private key of a candidate.
func DeclarePrivateKey(c *gin.Context) {
	username, err := utils.GetSessionID(c)
	if err != nil {
		c.String(http.StatusForbidden, "You need to be logged in.")
		return
	}

	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusNotFound, "Candidate not found.")
		return
	}

	if candidate.KeyState != config.KeysGenerated {
		c.String(http.StatusBadRequest, "Candidate has already declared private key.")
		return
	}

	keys := declarePrivateKey_Keys{}
	err = c.BindJSON(&keys)
	if err != nil {
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	candidate.PublicKey = keys.PublicKey
	candidate.PrivateKey = keys.PrivateKey
	candidate.KeyState = config.KeysDeclared
	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while declaring private keys of candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	c.JSON(http.StatusOK, "Public Key succesfully received.")
}

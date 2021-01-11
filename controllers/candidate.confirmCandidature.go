package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/utils"
)

// Struct to accept new keys from the client.
type confirmCandidature_Keys struct {
	PublicKey  string `json:"pubkey"`
	PrivateKey string `json:"privkey"`
}

// API handler to confirm candidature by updating the public and private keys of a candidate.
func ConfirmCandidature(c *gin.Context) {
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

	if candidate.KeyState != config.KeysNotGenerated {
		c.String(http.StatusBadRequest, "Candidate has already registered.")
		return
	}

	keys := confirmCandidature_Keys{}
	err = c.BindJSON(&keys)
	if err != nil {
		c.String(http.StatusBadRequest, "Data format not recognized.")
		return
	}

	candidate.PublicKey = keys.PublicKey
	candidate.PrivateKey = keys.PrivateKey
	candidate.KeyState = config.KeysGenerated
	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while confirming candidature: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database Error.")
		return
	}

	c.JSON(http.StatusOK, "Candidature confirmed.")
}

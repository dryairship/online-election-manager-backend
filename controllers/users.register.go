package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterNewVoter(c *gin.Context) {
	roll := c.PostForm("roll")
	passHash := c.PostForm("pass")
	authCode := c.PostForm("auth")

	if roll == "CEO" {
		RegisterCEO(c)
		return
	}

	if roll[0] == 'P' {
		RegisterCandidate(c)
		return
	}

	voter, err := ElectionDb.FindVoter(roll)
	if err != nil {
		c.String(http.StatusForbidden, "You need to get a verification mail before you register.")
		return
	}

	if voter.AuthCode == "" {
		log.Println("[WARN] Re-registration attempt: ", voter)
		c.String(http.StatusForbidden, "Student has already registered.")
		return
	}

	if voter.AuthCode != authCode {
		c.String(http.StatusBadRequest, "Wrong authentication code.")
		return
	}

	if passHash == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		log.Println("[WARN] Registration attempt with empty password: ", voter)
		c.String(http.StatusBadRequest, "Empty password is not allowed.")
		return
	}

	voter.Password = passHash
	voter.AuthCode = ""

	err = ElectionDb.UpdateVoter(roll, &voter)
	if err != nil {
		log.Println("[ERROR] Database error while registering voter: ", voter)
		c.String(http.StatusInternalServerError, "Database Error")
	} else {
		c.String(http.StatusAccepted, "Voter successfully registered.")
	}
}

// API handler to register a new candidate.
func RegisterCandidate(c *gin.Context) {
	username := c.PostForm("roll")
	passHash := c.PostForm("pass")
	authCode := c.PostForm("auth")

	candidate, err := ElectionDb.GetCandidate(username)

	if err != nil {
		c.String(http.StatusNotFound, "Candidate not found.")
		return
	}

	if candidate.Password != "" {
		c.String(http.StatusForbidden, "Candidate has already registered.")
		return
	} else if candidate.AuthCode == "" {
		c.String(http.StatusForbidden, "You need to get a verification<br>mail before you register.")
		return
	}

	if candidate.AuthCode != authCode {
		c.String(http.StatusBadRequest, "Wrong authentication code.")
		return
	}

	if passHash == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		c.String(http.StatusBadRequest, "Empty password is not allowed.")
		return
	}

	candidate.Password = passHash
	candidate.AuthCode = ""

	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while registering candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
	} else {
		c.String(http.StatusAccepted, "Candidate successfully registered.")
	}
}

// API handler to register CEO's account.
func RegisterCEO(c *gin.Context) {
	passHash := c.PostForm("pass")
	authCode := c.PostForm("auth")

	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		c.String(http.StatusForbidden, "You need to get a verification mail before you register.")
		return
	}

	if ceo.AuthCode == "" {
		log.Println("[WARN] CEO tried to re-register")
		c.String(http.StatusForbidden, "CEO has already registered.")
		return
	}

	if ceo.AuthCode != authCode {
		c.String(http.StatusBadRequest, "Wrong authentication code.")
		return
	}

	ceo.Password = passHash
	ceo.AuthCode = ""

	err = ElectionDb.UpdateCEO(&ceo)
	if err != nil {
		log.Println("[ERROR] Database error while registering CEO: ", ceo, err.Error())
		c.String(http.StatusInternalServerError, "Database Error")
	} else {
		c.String(http.StatusAccepted, "CEO successfully registered.")
	}
}

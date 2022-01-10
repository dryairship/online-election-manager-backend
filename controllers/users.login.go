package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/utils"
)

// API handler to check the user's login credentials.
func CheckUserLogin(c *gin.Context) {
	roll := c.PostForm("roll")
	passHash := c.PostForm("pass")
	ecPassHash := c.PostForm("ecpass")

	if config.UsingCaptcha {
		captchaId := c.PostForm("captchaId")
		captchaValue := c.PostForm("captchaValue")

		captchaSuccess := utils.VerifyCaptcha(captchaId, captchaValue)
		if !captchaSuccess {
			c.String(http.StatusBadRequest, "Incorrect CAPTCHA")
			return
		}
	}

	if roll == "CEO" {
		CEOLogin(c)
		return
	}

	if roll[0] == 'P' {
		CandidateLogin(c)
		return
	}

	if config.ElectionState == config.VotingNotYetStarted {
		c.String(http.StatusForbidden, "Voting has not yet started.")
		return
	}

	voter, err := ElectionDb.FindVoter(roll)
	if err != nil {
		c.String(http.StatusForbidden, "This student has not registered.")
		return
	}

	if !voter.AtHome {
		if ecPassHash != "" {
			if ecPassHash != config.CampusPassword {
				c.String(http.StatusForbidden, "Invalid Password.")
				return
			}
		} else {
			c.String(http.StatusForbidden, "Please provide the administrator password.")
			return
		}
	} else {
		if ecPassHash != config.DefaultAdminPassword {
			c.String(http.StatusForbidden, "Invalid admin Password.")
			return
		}
	}

	if voter.Password != passHash {
		c.String(http.StatusForbidden, "Invalid Password.")
		return
	}

	utils.StartSession(c)

	simplifiedVoter := voter.Simplify()
	c.JSON(http.StatusOK, &simplifiedVoter)
}

// API handler to check candidate's login credentials.
func CandidateLogin(c *gin.Context) {
	username := c.PostForm("roll")
	passHash := c.PostForm("pass")
	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusForbidden, "This candidate has not yet registered.")
		return
	}

	if candidate.Password != passHash {
		c.String(http.StatusForbidden, "Invalid Password.")
		return
	}

	utils.StartSession(c)

	simplifiedCandidate := candidate.Simplify()
	c.JSON(http.StatusOK, &simplifiedCandidate)
}

// API handler to check CEO's login credentials.
func CEOLogin(c *gin.Context) {
	passHash := c.PostForm("pass")
	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		log.Println("[ERROR] CEO Login attempted, but CEO has not registered")
		c.String(http.StatusForbidden, "CEO has not yet registered.")
		return
	}

	if ceo.Password != passHash {
		c.String(http.StatusForbidden, "Invalid Password.")
		return
	}

	utils.StartSession(c)

	c.JSON(http.StatusOK, &ceo)
}

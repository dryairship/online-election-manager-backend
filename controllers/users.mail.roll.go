package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/utils"
)

type UserError struct {
	reason string
}

func (err *UserError) Error() string {
	return err.reason
}

// Function to check if a verification mail can be sent to the student.
func canMailBeSentToStudent(roll string) (bool, error) {
	voter, err := ElectionDb.FindVoter(roll)
	if err != nil {
		return true, nil
	} else {
		if voter.AuthCode == "" {
			return false, &UserError{"Student has already registered."}
		} else {
			return false, &UserError{"Verification mail has already been sent to this student."}
		}
	}
}

// API handler to send a verification mail to the student.s
func SendMailToStudent(c *gin.Context) {
	if config.ElectionState == config.VotingStopped {
		log.Println("[WARN] Verification-mail attempt after end of registration period: ", c.Param("roll"))
		c.String(http.StatusForbidden, "Registration period is over.")
		return
	}

	var captcha CAPTCHA
	err := c.BindJSON(&captcha)
	if err != nil {
		log.Println("[ERROR] CAPTCHA data failed to bind to struct: ", err.Error())
		c.String(http.StatusBadRequest, "Invalid data format")
		return
	}

	success := utils.VerifyCaptcha(captcha.Id, captcha.Value)
	if !success {
		c.String(http.StatusBadRequest, "Incorrect CAPTCHA")
		return
	}

	roll := c.Param("roll")

	if roll == "CEO" {
		SendMailToCEO(c)
		return
	}

	if roll[0] == 'P' {
		SendMailToCandidate(c)
		return
	}

	_, err = canMailBeSentToStudent(roll)
	if err != nil {
		log.Println("[WARN] Verification-mail not sent to voter: ", roll, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	skeleton, err := ElectionDb.FindStudentSkeleton(roll)
	if err != nil {
		log.Println("[ERROR] Student not found in database while sending verification mail: ", roll, err.Error())
		c.String(http.StatusNotFound, "Invalid Roll Number. Student does not exist.")
		return
	}

	voter := skeleton.CreateVoter(utils.GetRandomAuthCode())
	recipient := voter.GetMailRecipient()
	err = utils.SendMailTo(&recipient, "a voter")
	if err != nil {
		log.Println("[ERROR] Mailer not working while sending mail to voter: ", voter, err.Error())
		c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
		return
	} else {
		err = ElectionDb.AddNewVoter(&voter)
		if err != nil {
			log.Println("[ERROR] Database error when sending mail to voter: ", voter, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}
	}
	c.String(http.StatusAccepted, "Verification Mail successfully sent to "+voter.Email)
}

// API handler to send verification mail to a candidate.
func SendMailToCandidate(c *gin.Context) {
	if config.ElectionState != config.VotingNotYetStarted {
		c.String(http.StatusForbidden, "Registration period for<br>candidates is over.")
		return
	}

	username := c.Param("roll")
	candidate, err := ElectionDb.GetCandidate(username)
	if err != nil {
		c.String(http.StatusNotFound, "This candidate is not<br>contesting elections.")
		return
	}

	if candidate.Password != "" {
		c.String(http.StatusForbidden, "This candidate has already<br>registered.")
		return
	} else if candidate.AuthCode != "" {
		c.String(http.StatusForbidden, "Verification mail has already<br>been sent to this candidate.")
		return
	}

	candidate.AuthCode = utils.GetRandomAuthCode()
	recipient := candidate.GetMailRecipient()
	err = utils.SendMailTo(&recipient, "a candidate")
	if err != nil {
		log.Println("[ERROR] Mailer error while sending mail to candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
		return
	}

	err = ElectionDb.UpdateCandidate(username, &candidate)
	if err != nil {
		log.Println("[ERROR] Database error while sending mail to candidate: ", candidate, err.Error())
		c.String(http.StatusInternalServerError, "Database error.")
		return
	}

	c.String(http.StatusAccepted, "Verification Mail successfully sent<br>to "+candidate.Email)
}

// API handler to send verification mail to CEO.
func SendMailToCEO(c *gin.Context) {
	ceo, err := ElectionDb.GetCEO()

	// Checking for err == nil because ceo is inserted in the database when
	// the verification mail is sent to them,
	if err == nil {
		log.Println("[WARN] CEO re-requested verification mail")
		c.String(http.StatusForbidden, "Verification mail has already been sent to CEO.")
		return
	}

	skeleton, err := ElectionDb.FindStudentSkeleton(config.RollNumberOfCEO)
	if err != nil {
		log.Println("[ERROR] No CEO in the database, can't send verification mail")
		c.String(http.StatusInternalServerError, "No CEO assigned.")
		return
	}

	ceo = skeleton.CreateCEO(utils.GetRandomAuthCode())
	recipient := ceo.GetMailRecipient()
	err = utils.SendMailTo(&recipient, "the CEO")
	if err != nil {
		log.Println("[ERROR] Database error while sending mail to CEO: ", ceo, err.Error())
		c.String(http.StatusInternalServerError, "Mailer Utility is not working.")
		return
	} else {
		err = ElectionDb.InsertCEO(&ceo)
		if err != nil {
			log.Println("[ERROR] Database error while inserting CEO: ", ceo, err.Error())
			c.String(http.StatusInternalServerError, "Database error.")
			return
		}
	}
	c.String(http.StatusAccepted, "Verification Mail successfully sent to "+ceo.Email)
}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/utils"
)

func EnsureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.GetSessionID(c)
		if err != nil {
			log.Printf("[WARN] Session ID error <%s> at <%s>", err.Error(), c.FullPath())
			c.String(http.StatusForbidden, "You need to be logged in.")
			c.Abort()
		} else {
			c.Set("ID", id)
			c.Next()
		}
	}
}

func EnsureCEO() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("ID")
		if id != "CEO" {
			log.Printf("[WARN] Non-CEO <%s> attempted to access <%s>", id, c.FullPath())
			c.String(http.StatusForbidden, "Only the CEO can access this.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureVotingNotYetStarted() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ElectionState != config.VotingNotYetStarted {
			c.String(http.StatusForbidden, "This should be done before voting is started.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureVotingStarted() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ElectionState != config.AcceptingVotes {
			c.String(http.StatusForbidden, "Not accepting votes currently.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureVotingStopped() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ElectionState != config.VotingStopped {
			c.String(http.StatusForbidden, "Voting has not been stopped.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureResultsCalculated() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ElectionState != config.ResultsCalculated {
			c.String(http.StatusForbidden, "Results have not yet been calculated.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureAfterVotingStarted() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ElectionState == config.VotingNotYetStarted {
			c.String(http.StatusForbidden, "Voting has not yet started.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func EnsureBeforeVotingStopped() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.ElectionState == config.VotingStopped || config.ElectionState == config.ResultsCalculated {
			c.String(http.StatusForbidden, "Voting has been stopped.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

package controllers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/models"
)

// Function to match voter's roll number with the regular expression for a post
// to determine whether or not the voter is eligible to vote for a post.
func GetPostsForVoter(roll string) []models.VotablePost {
	votablePosts := []models.VotablePost{}
	for _, post := range Posts {
		pattern := post.VoterRegex
		canVote, err := regexp.MatchString(pattern, roll)
		if err == nil && canVote && !post.Resolved {
			votablePosts = append(votablePosts, post.ConvertToVotablePost())
		}
	}
	return votablePosts
}

// API handler to fetch all posts for which this voter is eligible to vote.
func GetVotablePosts(c *gin.Context) {
	roll := c.GetString("ID")
	votablePosts := GetPostsForVoter(roll)
	c.JSON(http.StatusOK, votablePosts)
}

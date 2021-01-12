package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function to match voter's roll number with the regular expression for a post
// to determine whether or not the voter is eligible to vote for a post.
func GetPostsForVoter(roll string) []string {
	voter, err := ElectionDb.FindVoter(roll)
	if err != nil {
		return nil
	}
	return voter.Posts
}

// API handler to fetch all posts for which this voter is eligible to vote.
func GetVotablePosts(c *gin.Context) {
	roll := c.GetString("ID")
	votablePosts := GetPostsForVoter(roll)
	c.JSON(http.StatusOK, votablePosts)
}

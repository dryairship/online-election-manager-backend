package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API handler to fetch all posts from the database.
func FetchPosts(c *gin.Context) {
	posts, err := ElectionDb.GetPostsForCEO()
	if err != nil {
		log.Println("[ERROR] Database error while fetching posts: ", err.Error())
		c.String(http.StatusInternalServerError, "Error while fetching posts.")
		return
	}
	c.JSON(http.StatusOK, &posts)
}

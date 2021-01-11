package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/utils"
)

func CEOMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.GetSessionID(c)
		if err != nil {
			log.Printf("[WARN] Session ID error <%s> at <%s>", err.Error(), c.FullPath())
			c.String(http.StatusForbidden, "Only the CEO can access this.")
		} else if id != "CEO" {
			log.Printf("[WARN] Non-CEO <%s> attempted to access <%s>", id, c.FullPath())
			c.String(http.StatusForbidden, "Only the CEO can access this.")
		} else {
			c.Next()
		}
	}
}

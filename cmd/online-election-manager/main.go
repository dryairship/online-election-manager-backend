package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	_ "github.com/dryairship/online-election-manager/controllers"
	"github.com/dryairship/online-election-manager/router"
	"github.com/dryairship/online-election-manager/utils"
)

func main() {
	sessionKey := utils.GetRandomString(32)
	sessionDb := cookie.NewStore([]byte(sessionKey))

	r := gin.Default()
	r.Use(sessions.Sessions("SessionData", sessionDb))
	router.SetUpRoutes(r)

	if err := r.Run(config.ApplicationPort); err != nil {
		log.Fatalln("[ERROR] Could not start the server: ", err.Error())
	}
}

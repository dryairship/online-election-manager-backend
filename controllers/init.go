package controllers

import (
	"log"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"github.com/dryairship/online-election-manager/models"
)

var ElectionDb db.ElectionDatabase
var Posts []models.Post

func init() {
	var err error
	ElectionDb, err = db.ConnectToDatabase()
	if err != nil {
		log.Fatalln("[ERROR] Could not establish database connection: ", err.Error())
	}

	Posts, err = ElectionDb.GetPosts()
	if err != nil {
		log.Fatalln("[ERROR] Could not set posts data:", err.Error())
	}

	ceo, err := ElectionDb.GetCEO()
	if err != nil {
		log.Println("[WARN] CEO has not registered.")
	} else {
		config.PublicKeyOfCEO = ceo.PublicKey
	}
}

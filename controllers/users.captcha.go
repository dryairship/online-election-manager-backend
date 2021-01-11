package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/utils"
)

type CAPTCHA struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

func GetCaptcha(c *gin.Context) {
	id, value := utils.CreateCaptcha()
	captcha := CAPTCHA{
		Id:    id,
		Value: value,
	}

	c.JSON(http.StatusOK, &captcha)
}

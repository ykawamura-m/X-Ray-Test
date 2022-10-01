package controller

import (
	"X-Ray-Test/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ユーザー情報の保存処理
// DynamoDBにユーザー情報を保存する
func Save(c *gin.Context) {
	paramID := c.PostForm("id")
	paramName := c.PostForm("name")
	paramEmail := c.PostForm("email")
	paramTel := c.PostForm("tel")
	paramDB := c.PostForm("db")

	db, err := strconv.Atoi(paramDB)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if paramID == "" {
		err := repository.RegisterUser(paramName, paramEmail, paramTel, db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		err = repository.UpdateUser(paramID, paramName, paramEmail, paramTel, db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	c.Redirect(http.StatusFound, "/")
}

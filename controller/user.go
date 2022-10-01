package controller

import (
	"X-Ray-Test/model"
	"X-Ray-Test/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	MODE_NEW  = "NEW"
	MODE_EDIT = "EDIT"
)

// ユーザー登録/編集画面
// ユーザー情報の入力フォームを表示する
func User(c *gin.Context) {
	var user model.User

	paramID := c.Query("id")
	paramDB := c.Query("db")

	mode := MODE_NEW
	if paramID != "" && paramDB != "" {
		mode = MODE_EDIT

		db, err := strconv.Atoi(paramDB)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		user, err = repository.GetUser(c, paramID, db)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}

	c.HTML(http.StatusOK, "user.html", gin.H{
		"user": user,
		"mode": mode,
	})
}

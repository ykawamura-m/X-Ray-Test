package controller

import (
	"X-Ray-Test/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// メイン画面
// 登録済のユーザー一覧を表示する
func Index(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"users": users,
	})
}

package controller

import (
	"X-Ray-Test/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ユーザー情報の削除処理
// DBから対象のユーザー情報を削除する
func Delete(c *gin.Context) {
	paramID := c.PostForm("id")
	paramDB := c.PostForm("db")

	db, err := strconv.Atoi(paramDB)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	err = repository.DeleteUser(paramID, db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Redirect(http.StatusFound, "/")
}

package main

import (
	"X-Ray-Test/controller"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PORT = 5000
)

// メイン関数
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("view/html/*.html")
	r.Static("/css", "view/css/")
	r.GET("/", controller.Index)
	r.GET("/user", controller.User)
	r.POST("/save", controller.Save)
	r.POST("/delete", controller.Delete)
	r.Run(":" + strconv.Itoa(PORT))
}

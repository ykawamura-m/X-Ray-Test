package main

import (
	"X-Ray-Test/controller"
	"X-Ray-Test/middleware"
	"strconv"

	"github.com/aws/aws-xray-sdk-go/awsplugins/beanstalk"
	"github.com/gin-gonic/gin"
)

const (
	PORT = 5000
)

// メイン関数
func main() {
	initXray()
	initGin()
}

// X-Ray初期化
func initXray() {
	beanstalk.Init()
}

// Gin初期化
func initGin() {
	r := gin.Default()
	r.LoadHTMLGlob("view/html/*.html")
	r.Static("/css", "view/css/")
	r.Use(middleware.XrayMiddleware())
	r.GET("/", controller.Index)
	r.GET("/user", controller.User)
	r.POST("/save", controller.Save)
	r.POST("/delete", controller.Delete)
	r.Run(":" + strconv.Itoa(PORT))
}

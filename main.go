package main

import (
	"main/database"
	"main/handler"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()

	router := gin.Default()

	//--------------------------Basic Routes
	router.GET("/articles", handler.Games())
	router.GET("/articles/:aid", handler.GameById())
	router.GET("/articles/:aid/comments", handler.Comments())

	//--------------------------Authentication Required Routes
	router.Use(middleware.CheckAuthorization())
	{
		router.POST("/articles/:aid/comments", handler.PostComment())
		router.DELETE("/articles/:aid/comments/:cid", handler.DeleteComment())
	}

	//--------------------------User Routes
	router.POST("/users/signup", handler.Signup())
	router.POST("/users/login", handler.Login())
	router.POST("/users/sudo", handler.Sudo())

	router.Run(":3693")
}

//env GOOS=linux GOARCH=arm go build -o api_arm

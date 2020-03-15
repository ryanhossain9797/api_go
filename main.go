package main

import (
	"context"
	"fmt"
	"log"
	"main/handler"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://zireael9797:hummerh2suv@cluster0-lbdsg.azure.mongodb.net/test?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database("wikiDB")

	router := gin.Default()

	router.GET("/articles", handler.Games(db))
	router.GET("/articles/:aid", handler.GameById(db))

	router.GET("/articles/:aid/comments", handler.Comments(db))
	router.POST("/articles/:aid/comments", handler.PostComment(db))
	router.DELETE("/articles/:aid/comments/:cid", handler.DeleteComment(db))

	router.POST("/users/signup", handler.Signup(db))
	router.POST("/users/login", handler.Login(db))
	router.POST("/users/sudo", handler.Sudo(db))

	router.Run(":3693")
}

//env GOOS=linux GOARCH=arm go build -o api_arm

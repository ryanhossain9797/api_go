package handler

import (
	"context"
	"fmt"
	"log"
	"main/database"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Comments() gin.HandlerFunc {
	return func(c *gin.Context) {
		articleId, _ := primitive.ObjectIDFromHex(c.Param("aid"))
		collection := database.DB.Collection("comments")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		cursor, err := collection.Find(ctx, models.Review{Article: articleId})
		if err != nil {
			log.Fatal(err)
		}
		results := make([]models.Review, 0)
		for cursor.Next(ctx) {
			var result models.Review
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, result)

		}
		c.JSON(http.StatusOK, gin.H{"count": len(results), "comments": results})
	}
}

func PostComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("POSTCOMMENT: called")
		gamecollection := database.DB.Collection("articles")
		articleId, _ := primitive.ObjectIDFromHex(c.Param("aid"))
		collection := database.DB.Collection("comments")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		gameres := gamecollection.FindOne(ctx, bson.M{"_id": articleId})

		if gameres.Err() != nil {
			fmt.Println("Aborting Review")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Game doesn't exist"})
		} else {
			fmt.Println("Adding Review")
			res, err := collection.InsertOne(ctx, models.Review{Article: articleId, Comment: c.PostForm("comment"), Username: c.PostForm("username")})
			if err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, err)
			} else {
				fmt.Println(res)
				c.JSON(http.StatusCreated, res)
			}
		}
	}
}

func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		articleId, _ := primitive.ObjectIDFromHex(c.Param("aid"))
		commentId, _ := primitive.ObjectIDFromHex(c.Param("cid"))
		username := c.GetHeader("username")
		if len(username) > 0 {
			collection := database.DB.Collection("comments")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			res, err := collection.DeleteOne(ctx, models.Review{Article: articleId, Id: commentId, Username: username})
			if err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, err)
			} else {
				fmt.Println("DELETED_COMMENT: ")
				c.JSON(http.StatusOK, res)
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not set"})
		}
	}
}

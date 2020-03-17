package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"main/database"
	"main/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Games() gin.HandlerFunc {
	return func(c *gin.Context) {
		collection := database.DB.Collection("articles")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		results := make([]models.Game, 0)
		for cursor.Next(ctx) {
			var result models.Game
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, result)

		}
		c.JSON(200, gin.H{"count": len(results), "articles": results})
	}
}

func GameById() gin.HandlerFunc {
	return func(c *gin.Context) {
		collection := database.DB.Collection("articles")
		id, _ := primitive.ObjectIDFromHex(c.Param("aid"))
		article := models.Game{}
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := collection.FindOne(ctx, models.Game{Id: id}).Decode(&article)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
		}
		c.JSON(http.StatusOK, gin.H{"article": gin.H{
			"_id":     article.Id,
			"title":   article.Title,
			"imgurl":  article.Imgurl,
			"content": article.Content,
		}})
	}
}

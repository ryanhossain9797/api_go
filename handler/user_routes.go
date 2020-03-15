package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"main/models"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var hs = jwt.NewHS256([]byte("secret"))

func Signup(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("SIGNUP: called")
		usercollection := db.Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		email := c.PostForm("email")
		username := c.PostForm("username")
		password := c.PostForm("password")
		if len(email) > 0 && len(username) > 0 && len(password) > 0 {
			//--------------------------------IF ALL FIELDS PROVIDED
			userres := usercollection.FindOne(ctx, models.User{Email: email})
			if userres.Err() == nil {
				//--------------------------------IF USER ALREADY EXISTS
				fmt.Println("SIGNUP: Existing User - Aborting Signup")
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User already exists"})
			} else {
				//--------------------------------IF USER DOESN'T EXIST
				fmt.Println("SIGNUP: New User - Adding User")

				hashedpass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
				if err == nil {
					//--------------------------------IF PASSWORD HASHING IS SUCCESSFUL
					res, err := usercollection.InsertOne(ctx, models.User{Email: email, Username: username, Hashedpass: string(hashedpass)})
					if err != nil {
						//--------------------------------IF DB INSERTAION FAILS
						log.Fatal(err)
						c.JSON(http.StatusInternalServerError, err)
					} else {
						//--------------------------------IF EVERYTHING SUCCESSFUL
						fmt.Println("SIGNUP: User added to db")
						fmt.Println(res)
						c.JSON(http.StatusCreated, gin.H{
							"status":   "success!",
							"email":    email,
							"username": username,
						})
					}
				} else {
					//--------------------------------IF PASSWORD HASHING FAILS
					log.Fatal(err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err,
					})
				}
			}
		} else {
			//--------------------------------IF ALL FIELDS NOT PROVIDED
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "You must provide username, email and password",
			})
		}
	}

}

type CustomPayload struct {
	jwt.Payload
}

func Login(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("LOGIN: called")
		usercollection := db.Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		email := c.PostForm("email")
		password := c.PostForm("password")

		if len(email) > 0 && len(password) > 0 {
			//--------------------------------IF ALL FIELDS PROVIDED
			userres := usercollection.FindOne(ctx, models.User{Email: email})
			if userres.Err() != nil {
				//--------------------------------IF USER DOESN'T EXIST
				fmt.Println("LOGIN: User Doesn't Exist - Aborting Login")
				c.JSON(http.StatusNotFound, gin.H{"error": "User doesn't exist"})
			} else {
				//--------------------------------IF USER EXISTS
				fmt.Println("LOGIN: User Exists - Attempting Login")
				user := models.User{}
				usererr := userres.Decode(&user)
				if usererr == nil {
					//--------------------------------IF USER OBJECT CREATED SUCCESSFULLY
					fmt.Println("LOGIN: User decoded")
					err := bcrypt.CompareHashAndPassword([]byte(user.Hashedpass), []byte(password))
					if err == nil {
						//--------------------------------IF PASSWORD MATCHES
						fmt.Println("LOGIN: Password matched")
						now := time.Now()
						pl := CustomPayload{
							Payload: jwt.Payload{
								Issuer:         "gbrlsnchs",
								Subject:        user.Id.Hex(),
								Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
								ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
								NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
								IssuedAt:       jwt.NumericDate(now),
								JWTID:          "foobar",
							},
						}

						token, err := jwt.Sign(pl, hs)
						if err != nil {
							//--------------------------------IF JWT SIGNING FAILS
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": err,
							})
						} else {
							//--------------------------------IF JWT SIGNING SUCCEEDES
							fmt.Println(gin.H{
								"token": string(token),
								"uid":   user.Id.Hex(),
							})
							c.JSON(http.StatusOK, gin.H{
								"token": string(token),
							})
						}

					} else {
						//--------------------------------IF PASSWORD DEOSN'T MATCH
						fmt.Println("LOGIN: Password Mismatch")
						c.JSON(http.StatusUnauthorized, gin.H{
							"error": "password mismatch",
						})
					}
				} else {
					//--------------------------------IF USER OBJECT DECODE FAILS
					fmt.Println("LOGIN: User Decode Failed")
					log.Fatal(usererr)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": usererr,
					})
				}
			}
		} else {
			//--------------------------------IF ALL FIELDS NOT PROVIDED
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "You must provide email and password",
			})
		}
	}
}

func Sudo(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("SUDO: called")
		token := []byte(c.GetHeader("token"))
		fmt.Print(string(token))
		if len(token) > 0 {
			//--------------------------------IF TOKEN PROVIDED
			var pl CustomPayload
			hd, err := jwt.Verify(token, hs, &pl)
			if err == nil {

				usercollection := db.Collection("users")
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				userId, _ := primitive.ObjectIDFromHex(pl.Subject)
				userres := usercollection.FindOne(ctx, bson.M{"_id": userId})
				if userres.Err() != nil {
					fmt.Println("Invalid User")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Game doesn't exist"})
				} else {
					user := models.User{}
					userres.Decode(&user)
					c.JSON(http.StatusAccepted, gin.H{
						"username": user.Username,
						"email":    user.Email,
						"header":   hd,
					})
				}
			} else {
				c.JSON(http.StatusForbidden, gin.H{
					"error": err,
				})
			}
		} else {
			//--------------------------------IF TOKEN NOT PROVIDED
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You must provide an authorization token",
			})
		}
	}
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	Id      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string             `json:"title,omitempty" bson:"title,omitempty" required:"true"`
	Content string             `json:"content,omitempty" bson:"content,omitempty" required:"true"`
	Imgurl  string             `json:"imgurl,omitempty" bson:"imgurl,omitempty" required:"true"`
}

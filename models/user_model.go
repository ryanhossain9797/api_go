package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty" required:"true"`
	Username   string             `json:"username,omitempty" bson:"username,omitempty" required:"true"`
	Hashedpass string             `json:"hashedpass,omitempty" bson:"hashedpass,omitempty" required:"true"`
}

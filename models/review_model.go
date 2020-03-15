package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty" required:"true"`
	Comment  string             `json:"comment,omitempty" bson:"comment,omitempty" required:"true"`
	Article  primitive.ObjectID `json:"article,omitempty" bson:"article,omitempty" required:"true"`
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Role     string             `json:"role" bson:"role"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

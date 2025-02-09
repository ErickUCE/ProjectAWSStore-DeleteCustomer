package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer representa el esquema de la colección en MongoDB
type Customer struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone_number" bson:"phone_number"`
	Address   string             `json:"address" bson:"address"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"required,email"`
	Phone        *string            `json:"phone" validate:"required"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken"`
	UserType     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserId       string             `json:"user_id"`
}

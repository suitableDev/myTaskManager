package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email                    *string            `bson:"email" json:"email" validate:"email,required"`
	Username                 *string            `bson:"username" json:"username" validate:"required,min=2,max=30"`
	Password                 *string            `bson:"password" json:"password,omitempty" validate:"required,min=6,max=30"`
	UserType                 *string            `bson:"user_type" json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	UserID                   *string            `bson:"user_id" json:"user_id,omitempty"`
	Token                    *string            `bson:"token" json:"token,omitempty"`
	RefreshToken             *string            `bson:"refresh_token" json:"refresh_token,omitempty"`
	CreatedAt                time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time          `bson:"updated_at" json:"updated_at"`
	VerificationToken        *string            `bson:"verification_token" json:"verification_token,omitempty"`
	Verified                 bool               `bson:"verified" json:"verified"`
	ResetPasswordToken       *string            `bson:"reset_password_token,omitempty" json:"reset_password_token,omitempty"`
	ResetPasswordTokenExpiry time.Time          `bson:"reset_password_token_expiry" json:"reset_password_token_expiry"`
}

type SignedDetails struct {
	Email    string
	Username string
	Uid      string
	UserType string
	jwt.StandardClaims
}

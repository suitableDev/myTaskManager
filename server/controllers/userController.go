package controllers

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"

	"task-manager/server/database"
)

var userCollection *mongo.Collection = database.GetUserCollection()
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup()

func Login()

func GetUsers()

func GetUser()

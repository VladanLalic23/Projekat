package dto

import (
	"projekat/backend/userService/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdministratorRegistrationDTO struct {
	Name        string              `json:"name"`
	Surname     string              `json:"surname"`
	Username    string              `json:"username"`
	Password    string              `json:"password"`
	Email       string              `json:"email"`
	PhoneNumber string              `json:"phoneNumber"`
	Gender      *model.Gender       `json:"gender"`
	BirthDate   *primitive.DateTime `json:"birthDate"`
}

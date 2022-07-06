package dto

import (
	"projekat/backend/userService/model"
)

type VerificationRequestDTO struct {
	UserId           string         `json:"_id"`
	Name             string         `json:"name"`
	Surname          string         `json:"surname"`
	ImageUrl         string         `json:"imageUrl"`
	VerificationType model.UserType `json:"verificationType"`
}

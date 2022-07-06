package dto

import (
	"projekat/backend/userService/model"
)

type VerificationDTO struct {
	Id               string         `json:"_id"`
	UserId           string         `json:"userId"`
	Name             string         `json:"verificationName"`
	Surname          string         `json:"verificationSurname"`
	ImageUrl         string         `json:"verificationImageUrl"`
	VerificationType model.UserType `json:"verificationType"`
}

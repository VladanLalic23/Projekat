package dto

import "projekat/backend/userService/model"

type UserVerificationDTO struct {
	UserId           string         `json:"_id"`
	VerificationType model.UserType `json:"verificationType"`
}

package dto

import (
	"projekat/backend/userService/model"
)

type RegularUserPostDTO struct {
	Id          string             `bson:"_id"`
	PrivacyType *model.PrivacyType `bson:"privacyType"`
}

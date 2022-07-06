package dto

import (
	"projekat/backend/userService/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegularUserProfileDataDTO struct {
	Id             primitive.ObjectID   `json:"_id"`
	Name           string               `json:"name"`
	Surname        string               `json:"surname"`
	Username       string               `json:"username"`
	Biography      string               `json:"biography"`
	WebSite        string               `json:"webSite"`
	ProfilePrivacy model.ProfilePrivacy `json:"profilePrivacy"`
}

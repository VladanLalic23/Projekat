package repository

import (
	"context"
	"fmt"
	"log"
	"projekat/backend/userService/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VerificationRequestRepository struct {
	Database *mongo.Database
}

func (repository *VerificationRequestRepository) CreateVerificationRequest(verificationRequest *model.VerificationRequest) error {
	verificationRequestCollection := repository.Database.Collection("verificationRequests")
	_, err := verificationRequestCollection.InsertOne(context.TODO(), &verificationRequest)
	if err != nil {
		return fmt.Errorf("verification request is NOT created")
	}
	return nil
}

func (repository *VerificationRequestRepository) GetAllVerificationRequests() ([]bson.D, error) {

	verificationRequestCollection := repository.Database.Collection("verificationRequests")
	filterCursor, err := verificationRequestCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var verificationRequestFiltered []bson.D
	if err = filterCursor.All(context.TODO(), &verificationRequestFiltered); err != nil {
		log.Fatal(err)
	}
	return verificationRequestFiltered, nil
}

func (repository *VerificationRequestRepository) DeleteVerificationRequest(id primitive.ObjectID) error {

	verificationRequestCollection := repository.Database.Collection("verificationRequests")
	_, err := verificationRequestCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

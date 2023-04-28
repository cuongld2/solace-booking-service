package mongodb

import (
	"context"
	"errors"

	"cuongld2.com/api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// meetingModel represent a mgo database session with a meeting data model
type MeetingModel struct {
	C *mongo.Collection
}

// All method will be used to get all records from meetings table
func (m *MeetingModel) All() ([]models.Meeting, error) {
	// Define variables
	ctx := context.TODO()
	b := []models.Meeting{}

	// Find all meetings
	meetingCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = meetingCursor.All(ctx, &b)
	if err != nil {
		return nil, err
	}

	return b, err
}

// FindByID will be used to find a meeting registry by id
func (m *MeetingModel) FindByID(id string) (*models.Meeting, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find meeting by id
	var meeting = models.Meeting{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&meeting)
	if err != nil {
		// Checks if the meeting was not found
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &meeting, nil
}

// Insert will be used to insert a new meeting registry
func (m *MeetingModel) Insert(meeting models.Meeting) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), meeting)
}

// Delete will be used to delete a meeting registry
func (m *MeetingModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// meeting is used to represent meeting profile data
type Meeting struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	SenderUserID       string             `bson:"sender_userid"`
	ReceiverUserID     []string           `bson:"receiver_userid"`
	MeetingTitle       string             `bson:"meeting_title"`
	MeetingDescription []string           `bson:"meeting_description"`
	MeetingTime        time.Time          `bson:"meeting_time"`
	CreatedTime        time.Time          `bson:"created_time"`
}

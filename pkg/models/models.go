package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// meeting is used to represent meeting profile data
type Meeting struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	SenderUserID       string             `json:"sender_userid"`
	ReceiverUserID     []string           `json:"receiver_userid"`
	MeetingTitle       string             `json:"meeting_title"`
	MeetingDescription string             `json:"meeting_description"`
	MeetingCategory    string             `json:"meeting_category"`
	MeetingTime        time.Time          `bson:"meeting_time"`
	CreatedTime        time.Time          `bson:"created_time"`
}

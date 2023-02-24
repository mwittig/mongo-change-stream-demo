package app

import "time"

type Event struct {
	MessageID string    `bson:"message_id"`
	Message   string    `bson:"message"`
	CreatedAt time.Time `json:"" bson:"created_at"`
}

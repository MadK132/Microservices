package mongo

import "time"

type Event struct {
	Source    string    `bson:"source"`
	Action    string    `bson:"action"`
	Timestamp time.Time `bson:"timestamp"`
}

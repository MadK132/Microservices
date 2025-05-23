// internal/entity/event.go
package entity

import "time"

type Event struct {
	Source    string    
	Action    string    
	Timestamp time.Time 
}

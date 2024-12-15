package models

import "time"

// Log represents a log entry in the system
type Log struct {
	ID        int       `json:"id"`
	EventTime time.Time `json:"event_time"`
	EventType string    `json:"event_type"`
	UserID    int       `json:"user_id"`
	Message   string    `json:"message"`
}

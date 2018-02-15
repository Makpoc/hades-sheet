package models

import "time"

// UserTime contains information about users and their time
type UserTime struct {
	UserName    string
	CurrentTime time.Time
	Offset      string
}

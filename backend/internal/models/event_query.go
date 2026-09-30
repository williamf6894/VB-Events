package models

import "time"

type EventQuery struct {
	Search string
	Before *time.Time
	After  *time.Time
	Full   *bool
}

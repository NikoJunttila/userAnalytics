// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Name        string
	Url         string
	TotalVisits int32
	TotalUnique int32
	TotalTime   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DomainFollow struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UserID     uuid.UUID
	DomainID   uuid.UUID
	DomainName string
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Passhash  string
	ApiKey    string
}

type Visit struct {
	Createdat     time.Time
	Visitorstatus string
	Visitduration int32
	Domain        uuid.UUID
	Visitfrom     string
}

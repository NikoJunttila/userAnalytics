package main

import (
	"time"

	"github.com/nikojunttila/userAnalytics/internal/database"
)

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}
type LoginUser struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	Email     string    `json:"email"`
}
type CurretUser struct {
	Name   string `json:"name"`
	APIKey string `json:"api_key"`
	Email  string `json:"email"`
}
type DomainCreate struct {
	Name     string `json:"name"`
	DomainId string `json:"domainID"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
		Email:     dbUser.Email,
		Password:  dbUser.Passhash,
	}
}
func databaseUserToLogin(dbUser database.User) LoginUser {
	return LoginUser{
		CreatedAt: dbUser.CreatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
		Email:     dbUser.Email,
	}
}
func databaseCurrentUser(dbUser database.User) CurretUser {
	return CurretUser{
		Name:   dbUser.Name,
		APIKey: dbUser.ApiKey,
		Email:  dbUser.Email,
	}
}

func databaseDomainToDomain(dbDomain database.Domain) Domain {
	return Domain{
		ID:          dbDomain.ID,
		OwnerID:     dbDomain.OwnerID,
		Name:        dbDomain.Name,
		Url:         dbDomain.Url,
		TotalVisits: dbDomain.TotalVisits,
		TotalUnique: dbDomain.TotalUnique,
		TotalTime:   dbDomain.TotalTime,
		CreatedAt:   dbDomain.CreatedAt,
		UpdatedAt:   dbDomain.UpdatedAt,
	}
}

type Domain struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	TotalVisits int64     `json:"total_visits"`
	TotalUnique int64     `json:"total_unique"`
	TotalTime   int64     `json:"total_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

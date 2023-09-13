package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

type User struct {
  ID uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Name string `json:"name"`
  APIKey string `json:"api_key"`
  Email string `json:"email"`
  Password string `json:"password"`
}
type LoginUser struct {
  CreatedAt time.Time `json:"created_at"`
  Name string `json:"name"`
  APIKey string `json:"api_key"`
  Email string `json:"email"`
}
func databaseUserToUser(dbUser database.User) User{
  return User{
    ID: dbUser.ID,
    CreatedAt: dbUser.CreatedAt,
    UpdatedAt: dbUser.UpdatedAt,
    Name: dbUser.Name,
    APIKey: dbUser.ApiKey,
    Email: dbUser.Email,
    Password: dbUser.Passhash,
  }
}
func databaseUserToLogin(dbUser database.User) LoginUser{
  return LoginUser{
    CreatedAt: dbUser.CreatedAt,
    Name: dbUser.Name,
    APIKey: dbUser.ApiKey,
    Email: dbUser.Email,
  }
}

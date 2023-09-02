package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/auth"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
  type parameters struct{
    Name string `json:"name"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  user,err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
    ID: uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Name: params.Name,
  })
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
  }
  
  respondWithJson(w, 200 , databaseUserToUser(user))
} 
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request){
  apiKey, err := auth.GetAPIKey(r.Header)
  if err != nil {
    respondWithError(w, 403, fmt.Sprintf("auth error: %v", err))
    return
  }
  user,err := apiCfg.DB.GetUserByAPIKey(r.Context(),apiKey)
  if err != nil {
    respondWithError(w, 400, fmt.Sprintf("couldnt get user: %v", err))
    return 
  }
  respondWithJson(w, 200, databaseUserToUser(user))
}


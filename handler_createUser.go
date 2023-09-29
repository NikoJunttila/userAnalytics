package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

  "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
  type parameters struct{
    Name string `json:"name"`
    Password string `json:"password"`
    Email string `json:"email"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  hashPassword := hashAndSalt([]byte(params.Password))
  user,err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
    ID: uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Name: params.Name,
    Email: params.Email,
    Passhash: hashPassword,
  })
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error with DB: %v", err))
  }
  
  respondWithJson(w, 200 , databaseUserToUser(user))
} 
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User){
  respondWithJson(w, 200, databaseCurrentUser(user))
}
func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
  type parameters struct{
    Email string `json:"email"`
    Password string `json:"password"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  user,err := apiCfg.DB.GetUserByEmail(r.Context(),params.Email)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error with DB: %v", err))
  }
  comparedPass := comparePasswords(user.Passhash, []byte(params.Password))
  if !comparedPass{
    respondWithError(w, 400, fmt.Sprint("Wrong password. Try again"))
    return
  }

  respondWithJson(w, 200 , databaseUserToLogin(user))
} 
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
    // Since we'll be getting the hashed password from the DB it
    // will be a string so we'll need to convert it to a byte slice
    byteHash := []byte(hashedPwd)
    err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
    if err != nil {
        fmt.Println(err)
        return false
    }
    return true
}
func hashAndSalt(pwd []byte) string {
    // Use GenerateFromPassword to hash & salt pwd.
    // MinCost is just an integer constant provided by the bcrypt
    // package along with DefaultCost & MaxCost. 
    // The cost can be any value you want provided it isn't lower
    // than the MinCost (4)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        fmt.Println(err)
    }
    // GenerateFromPassword returns a byte slice so we need to
    // convert the bytes to a string and return it
    return string(hash)
}

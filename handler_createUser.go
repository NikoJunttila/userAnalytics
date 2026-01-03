package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		//Name string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	hashPassword := hashAndSalt([]byte(params.Password))

	apiKeyBytes := make([]byte, 32)
	_, err = rand.Read(apiKeyBytes)
	if err != nil {
		respondWithError(w, 500, "could not generate api key")
		return
	}
	apiKey := hex.EncodeToString(apiKeyBytes)

	userID := uuid.New().String()
	err = apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "placeholder",
		ApiKey:    apiKey,
		Email:     params.Email,
		Passhash:  hashPassword,
	})

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			respondWithError(w, 400, "Email already taken")
			return
		}
		// Handle other database errors
		errMsg := fmt.Sprintf("error with DB: %v", err)
		fmt.Println(errMsg)
		respondWithError(w, 400, errMsg)
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error fetching created user: %v", err))
		return
	}

	respondWithJson(w, 200, databaseUserToUser(user))
}
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, databaseCurrentUser(user))
}
func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, "no users found with this email/password combo")
		return
	}
	comparedPass := comparePasswords(user.Passhash, []byte(params.Password))
	if !comparedPass {
		respondWithError(w, 400, "Wrong password. Try again")
		return
	}

	respondWithJson(w, 200, databaseUserToLogin(user))
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

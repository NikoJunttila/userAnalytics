package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
  "crypto/rand"
	"encoding/hex"
	"time"

	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerChangePass(w http.ResponseWriter, r *http.Request,user database.User){
  type parameters struct{
    OldPass string `json:"oldPass"`
    NewPass string `json:"newPass"`
  }
  decoder := json.NewDecoder(r.Body)
  params :=  parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  comparedPass := comparePasswords(user.Passhash, []byte(params.OldPass))
  if !comparedPass{
    respondWithError(w, 400, fmt.Sprint("Old password not correct"))
    return
  }
  hashPassword := hashAndSalt([]byte(params.NewPass))
  err = apiCfg.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
 Passhash: hashPassword,
 ID: user.ID,
  })
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error updating password: %v", err))
    return 
  }
 respondWithJson(w, 201, fmt.Sprint("Updated password")) 
}

func (apiCfg *apiConfig) handlerForgotPass(w http.ResponseWriter, r *http.Request, emailSecret string){
  type parameters struct{
    Email string `json:"email"`
  }
  decoder := json.NewDecoder(r.Body)
  params :=  parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  recipient := params.Email
  _, err = apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("No user found with this email: %v", err))
    return
  }
	// Generate a unique token
	token, err := generateToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set expiration time (e.g., 1 hour from now)
	// expirationTime := time.Now().Add(time.Hour)

	// Save the password reset information to new postGres table
	// resetInfo := PasswordReset{
	// 	Email:          recipient,
	// 	ResetToken:     token,
	// 	ExpirationTime: expirationTime,
	// }
  
  resetURL := fmt.Sprintf("http://localhost:800/v1/reset-password?token=%s", token)

  auth := smtp.PlainAuth("", "nikosamulijunttila@gmail.com", emailSecret, "smtp.gmail.com")
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{recipient}
	msg := []byte("To:" + recipient + "\r\n" +
		"Subject: Forgot password\r\n" +
		"\r\n" +
		"User analytics service password reset. \n To reset your password follow link. " + resetURL + " \nIf you didn't order password reset you can ignore this message.\r\n")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "Niko.Junttila@gmail.com", to, msg)
	if err != nil {
    fmt.Println(err)
    respondWithError(w,400,fmt.Sprintf("error sending reset email: %v", err))
		return
	}
 respondWithJson(w, 200, fmt.Sprint("Sent password reset link to your email.")) 
}
type PasswordReset struct {
	Email           string
	ResetToken      string
	ExpirationTime  time.Time
}
func generateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

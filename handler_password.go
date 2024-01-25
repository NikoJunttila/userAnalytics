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
  // user.ID
}

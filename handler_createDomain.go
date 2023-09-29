package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateDomain(w http.ResponseWriter, r *http.Request, user database.User){
  type parameters struct{
    Name string `json:"name"`
    Url string `json:"url"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  domainUUID := uuid.New()
  domain,err := apiCfg.DB.CreateDomain(r.Context(), database.CreateDomainParams{
    ID: domainUUID,
    OwnerID: user.ID,
    Name: params.Name,
    Url: params.Url,
    TotalVisits: 0,
    TotalUnique: 0,
    TotalTime: 0,
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
  })
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
  }
  _, err = apiCfg.DB.CreateDomainFollow(r.Context(), database.CreateDomainFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UserID:    user.ID,
		DomainID:  domainUUID,
    DomainName: params.Name,
	})
	if err != nil {
    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
		return
	}

  respondWithJson(w, 201 , domain)
}


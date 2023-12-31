package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (cfg *apiConfig) handlerDomainFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	domainFollows, err := cfg.DB.GetDomainsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}
	respondWithJson(w, 200, domainFollows)
}

func (cfg *apiConfig) handlerDomainFollowCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
    DomainID uuid.UUID `json:"domain_id"`
    DomainName string `json:"domain_name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	domainFollow, err := cfg.DB.CreateDomainFollow(r.Context(), database.CreateDomainFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UserID:    user.ID,
		DomainID:    params.DomainID,
    DomainName: params.DomainName, 
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	respondWithJson(w, 201, domainFollow)
}
func (apiCfg *apiConfig) handlerGetDomain(w http.ResponseWriter, r *http.Request, user database.User) {
		type parameters struct {
    DomainID uuid.UUID `json:"domain_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	stats, err := apiCfg.DB.GetTotalCount(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	err = apiCfg.DB.UpdateDomain(r.Context(), database.UpdateDomainParams{
		ID:          params.DomainID,
		TotalVisits: int32(stats.TotalCount),
		TotalUnique: int32(stats.NewVisitorCount),
		TotalTime:   int32(stats.AvgVisitDuration),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
  domain, err := apiCfg.DB.GetDomain(r.Context(), params.DomainID)
  if err != nil { 
    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting domain: %v",err))
    return
  }
	respondWithJson(w, 200, domain)
}

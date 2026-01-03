package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerDomainFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	domainFollows, err := apiCfg.DB.GetDomainsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=200")
	respondWithJson(w, 200, domainFollows)
}

func (apiCfg *apiConfig) handlerDomainFollowCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		DomainID string `json:"domain_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	domain, _ := apiCfg.DB.GetDomain(r.Context(), params.DomainID)
	domainFollowID := uuid.New().String()
	err = apiCfg.DB.CreateDomainFollow(r.Context(), database.CreateDomainFollowParams{
		ID:         domainFollowID,
		CreatedAt:  time.Now().UTC(),
		UserID:     user.ID,
		DomainID:   params.DomainID,
		DomainName: domain.Name,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	domainFollow := database.DomainFollow{
		ID:         domainFollowID,
		CreatedAt:  time.Now().UTC(),
		UserID:     user.ID,
		DomainID:   params.DomainID,
		DomainName: domain.Name,
	}

	respondWithJson(w, 201, domainFollow)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateDomain(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	domainID := uuid.New().String()
	err = apiCfg.DB.CreateDomain(r.Context(), database.CreateDomainParams{
		ID:          domainID,
		OwnerID:     user.ID,
		Name:        params.Name,
		Url:         params.Url,
		TotalVisits: 0,
		TotalUnique: 0,
		TotalTime:   0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error creating domain: %v", err))
		return
	}
	err = apiCfg.DB.CreateDomainFollow(r.Context(), database.CreateDomainFollowParams{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		UserID:     user.ID,
		DomainID:   domainID,
		DomainName: params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error following domain: %v", err))
		return
	}

	domain, err := apiCfg.DB.GetDomain(r.Context(), domainID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error fetching created domain: %v", err))
		return
	}

	respondWithJson(w, 201, domain)
}
func (apiCfg *apiConfig) handlerGetDomain(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		DomainID string `json:"domain_id"`
	}
	/*   type compare struct {
	     Total float64 `json:"total"`
	     Unique float64 `json:"unique"`
	     } */
	type extendDomain struct {
		database.Domain
		Total  float64 `json:"total"`
		Unique float64 `json:"unique"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	domainID := params.DomainID
	domain, err := apiCfg.DB.GetDomain(r.Context(), domainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting domain: %v", err))
		return
	}
	stats1, err := apiCfg.DB.GetMonthStats(r.Context(), domainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	stats2, err := apiCfg.DB.GetPrevMonthStats(r.Context(), domainID)
	w.Header().Set("Cache-Control", "public, max-age=100")
	if err != nil || stats2.TotalCount == 0 {
		var infinite extendDomain
		infinite.Domain = domain
		infinite.Total = 0.0
		infinite.Unique = 0.0
		respondWithJson(w, 200, infinite)
		return
	}
	var stats extendDomain
	stats.Total = percentageDiff(stats1.TotalCount, stats2.TotalCount)
	stats.Unique = percentageDiff(stats1.NewVisitorCount, stats2.NewVisitorCount)
	stats.Domain = domain
	respondWithJson(w, 200, stats)

}
func percentageDiff(first int64, second int64) float64 {
	diff := (float64(first) - float64(second)) / float64(second) * 100.0
	return diff
}

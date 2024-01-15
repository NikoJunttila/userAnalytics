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
    return
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
func percentageDiff(first int64, second int64)float64{
  diff := (float64(first) - float64(second)) / float64(second) * 100.0
  return diff
}
func (apiCfg *apiConfig) handlerCompare(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		DomainID uuid.UUID `json:"domain_id"`
	}
  type compare struct {
  Total float64 `json:"total"`
  Unique float64 `json:"unique"`
  }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	stats1, err := apiCfg.DB.GetMonthStats(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	stats2, err := apiCfg.DB.GetPrevMonthStats(r.Context(), params.DomainID)
	if err != nil || stats2.TotalCount == 0 {
  var infinite compare
  infinite.Total = 0.0
  infinite.Unique = 0.0
  respondWithJson(w, 200, infinite)
  return
  }
  var stats compare;
  stats.Total = percentageDiff(stats1.TotalCount, stats2.TotalCount)
  stats.Unique = percentageDiff(stats1.NewVisitorCount, stats2.NewVisitorCount)
  w.Header().Set("Cache-Control", "public, max-age=3600")
  respondWithJson(w, 200, stats)
}

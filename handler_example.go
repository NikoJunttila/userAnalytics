package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)
func (apiCfg *apiConfig) handlerGetFreeDomain(w http.ResponseWriter, r *http.Request) {  
  type compare struct {
  Total float64 `json:"total"`
  Unique float64 `json:"unique"`
  }
  type extendDomain struct {
  database.Domain
  Total float64 `json:"total"`
  Unique float64 `json:"unique"`
  }
  domainIDString := chi.URLParam(r, "id")
  domainID, err := uuid.Parse(domainIDString)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return
	}
  domain, err := apiCfg.DB.GetDomain(r.Context(), domainID)
  if err != nil { 
    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting domain: %v",err))
    return
  }
  stats1, err := apiCfg.DB.GetMonthStats(r.Context(), domainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
  w.Header().Set("Cache-Control", "public, max-age=100")
	stats2, err := apiCfg.DB.GetPrevMonthStats(r.Context(), domainID)
	if err != nil || stats2.TotalCount == 0 {
  var infinite extendDomain
  infinite.Domain = domain
  infinite.Total = 0.0
  infinite.Unique = 0.0
  respondWithJson(w, 200, infinite)
  return
  }
  var stats extendDomain;
  stats.Total = percentageDiff(stats1.TotalCount, stats2.TotalCount)
  stats.Unique = percentageDiff(stats1.NewVisitorCount, stats2.NewVisitorCount)
  stats.Domain = domain
	respondWithJson(w, 200, stats)
}

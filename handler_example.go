
package main

import (
	"net/http"
  "fmt"
	"github.com/google/uuid"
)
func (apiCfg *apiConfig) handlerFreeCompare(w http.ResponseWriter, r *http.Request) {
  type compare struct {
  Total float64 `json:"total"`
  Unique float64 `json:"unique"`
  }

  //test below
  // domainIDString := "9c698b28-4a0c-49e2-815d-0ab446088352"
  domainIDString := "93417e06-8dc7-40ed-a9a5-d65a72fc5098"
  domainID, err := uuid.Parse(domainIDString)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return
	}
	stats1, err := apiCfg.DB.GetMonthStats(r.Context(), domainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	stats2, err := apiCfg.DB.GetPrevMonthStats(r.Context(), domainID)
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

func (apiCfg *apiConfig) handlerGetFreeDomain(w http.ResponseWriter, r *http.Request) {  
  //test below
  // domainIDString := "9c698b28-4a0c-49e2-815d-0ab446088352"
  domainIDString := "93417e06-8dc7-40ed-a9a5-d65a72fc5098"
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
  w.Header().Set("Cache-Control", "public, max-age=3600")
	respondWithJson(w, 200, domain)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateVisit(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Credentials", "true") 
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	type parameters struct {
		VisitStat     string    `json:"status"`
		VisitDuration int32     `json:"visitDuration"`
		Domain        uuid.UUID `json:"domain"`
		VisitFrom     string    `json:"visitFrom"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("error: %v \n", err)
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	dbCtx := context.Background()
	// Asynchronously save the visit to the database
	go func() {
		_, err := apiCfg.DB.CreateVisit(dbCtx, database.CreateVisitParams{
			Createdat:     time.Now().UTC(),
			Visitorstatus: params.VisitStat,
			Visitduration: params.VisitDuration,
			Domain:        params.Domain,
			Visitfrom:     params.VisitFrom,
		})
		if err != nil {
			fmt.Printf("error: %v \n", err)
		}
	}()

  var uniqueVisit int32 = 0
  if params.VisitStat == "new"{
  uniqueVisit = 1}
	err = apiCfg.DB.UpdateDomain(r.Context(), database.UpdateDomainParams{
		  ID:          params.Domain,
		  TotalVisits: 1,
		  TotalUnique: uniqueVisit,
	})
  if err != nil {
		fmt.Printf("error: %v \n", err)
    return 
  }
  fmt.Println("new  visit")
	respondWithJson(w, 200, "success")
}

func (apiCfg *apiConfig) handlerLimitedVisits(w http.ResponseWriter, r *http.Request) {
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
	stats, err := apiCfg.DB.GetLimitedCount(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	respondWithJson(w, 200, stats)
}

// I CBA with this
func (apiCfg *apiConfig) handlerSevenVisits(w http.ResponseWriter, r *http.Request) {
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
	stats, err := apiCfg.DB.GetSevenDays(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	respondWithJson(w, 200, stats)
}
func (apiCfg *apiConfig) handlerNinetyVisits(w http.ResponseWriter, r *http.Request) {
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
	stats, err := apiCfg.DB.GetNinetyDays(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	respondWithJson(w, 200, stats)
}

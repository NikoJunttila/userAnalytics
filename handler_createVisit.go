package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
  "context"
	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateVisit(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		VisitStat      string    `json:"status"`
		VisitDuration  int32     `json:"visitDuration"`
		Domain         uuid.UUID `json:"domain"`
		VisitFrom      string    `json:"visitFrom"`
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
		 if err != nil {
		     fmt.Printf("Database error: %v\n", err)
		}
	}()

	// Respond to the HTTP request immediately, without waiting for the database operation.
	respondWithJson(w, 200, nil)
}

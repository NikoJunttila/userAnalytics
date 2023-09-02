package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

func (apiCfg *apiConfig) handlerCreateVisit(w http.ResponseWriter, r *http.Request){
  type parameters struct{
    ID string `json:"ID"`
    Country string `json:"country"`
    IP string `json:"ip"`
    VisitStat string `json:"status"`
    Domain uuid.UUID `json:"domain"`
    VisitFrom string `json:"visitFrom"`
  }
  decoder := json.NewDecoder(r.Body)
  params := parameters{}
  err := decoder.Decode(&params)
  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
    return
  }
  visit,err := apiCfg.DB.CreateVisit(r.Context(), database.CreateVisitParams{
    ID: params.ID,
    Createdat: time.Now().UTC(),
    Country: params.Country,
    Ip: params.IP,
    Visitorstatus: params.VisitStat,
    Domain: params.Domain,
    Visitfrom: params.VisitFrom,
  })

  if err != nil {
    respondWithError(w,400,fmt.Sprintf("error parsing JSON: %v", err))
  }

  respondWithJson(w, 200 , visit)
}

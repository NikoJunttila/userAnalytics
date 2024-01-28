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
//NOBODY LOOK HERE SECRET STUFF AND NOT SUPER BAD CODE
func (apiCfg *apiConfig) handlerCreateVisit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	type parameters struct {
		VisitStat     string    `json:"status"`
		VisitDuration int32     `json:"visitDuration"`
		Domain        uuid.UUID `json:"domain"`
		VisitFrom     string    `json:"visitFrom"`
    Browser       string    `json:"browser"`
    Device        string    `json:"device"`
    OS            string    `json:"os"`
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
      Device:        params.Device,
      Os:            params.OS,
      Browser:       params.Browser,
		})
		if err != nil {
			fmt.Printf("error: %v \n", err)
		}
  var uniqueVisit int32 = 0
  if params.VisitStat == "new"{
  uniqueVisit = 1}
	err = apiCfg.DB.UpdateDomain(dbCtx, database.UpdateDomainParams{
		  ID:          params.Domain,
		  TotalVisits: 1,
		  TotalUnique: uniqueVisit,
	})
  if err != nil {
		fmt.Printf("error: %v \n", err)
  }
	}()
	respondWithJson(w, 200, "success")
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
  os , err := apiCfg.DB.GetOsCount7(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
  device , err := apiCfg.DB.GetDeviceCount7(r.Context(), params.DomainID)
 	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
  browser , err := apiCfg.DB.GetBrowserCount7(r.Context(), params.DomainID)
  if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
  type response struct {
    Original []database.GetSevenDaysRow `json:"original"`
    Os  []database.GetOsCount7Row `json:"os"`
    Device  []database.GetDeviceCount7Row `json:"device"`
    Browser  []database.GetBrowserCount7Row `json:"browser"`
  }
  	resp := response{
		Original: stats,
		Os:   os,
    Device: device,
    Browser: browser,
	}
	respondWithJson(w, 200, resp)
}
//30 days
func (apiCfg *apiConfig) handlerLimitedVisits(w http.ResponseWriter, r *http.Request) {
startTime := time.Now()
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
os, err := apiCfg.DB.GetOsCount30(r.Context(), params.DomainID)
if err != nil {
	respondWithError(w, http.StatusInternalServerError, "DB error for Os")
	return
}

device, err := apiCfg.DB.GetDeviceCount30(r.Context(), params.DomainID)
if err != nil {
	respondWithError(w, http.StatusInternalServerError, "DB error for Device")
	return
}

browser, err := apiCfg.DB.GetBrowserCount30(r.Context(), params.DomainID)
if err != nil {
	respondWithError(w, http.StatusInternalServerError, "DB error for Browser")
	return
}

type response30 struct {
	Original []database.GetLimitedCountRow `json:"original"`
	Os       []database.GetOsCount30Row        `json:"os"`
	Device   []database.GetDeviceCount30Row    `json:"device"`
	Browser  []database.GetBrowserCount30Row   `json:"browser"`
}

resp := response30{
	Original: stats,
	Os:       os,
	Device:   device,
	Browser:  browser,
}	
	elapsed := time.Since(startTime)
	fmt.Printf("my normal func took %s\n", elapsed)
respondWithJson(w, 200,resp)
}


func (apiCfg *apiConfig) handlerNinetyVisits(w http.ResponseWriter, r *http.Request) {
startTime := time.Now()
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

	// Create channels to receive results from goroutines
	statsChan := make(chan []database.GetNinetyDaysRow)
	osChan := make(chan []database.GetOsCount90Row)
	deviceChan := make(chan []database.GetDeviceCount90Row)
	browserChan := make(chan []database.GetBrowserCount90Row)
	errChan := make(chan error)

	// Use goroutines to execute queries concurrently
	go func() {
		stats, err := apiCfg.DB.GetNinetyDays(r.Context(), params.DomainID)
		if err != nil {
			errChan <- err
			return
		}
		statsChan <- stats
	}()

	go func() {
		os, err := apiCfg.DB.GetOsCount90(r.Context(), params.DomainID)
		if err != nil {
			errChan <- err
			return
		}
		osChan <- os
	}()

	go func() {
		device, err := apiCfg.DB.GetDeviceCount90(r.Context(), params.DomainID)
		if err != nil {
			errChan <- err
			return
		}
		deviceChan <- device
	}()

	go func() {
		browser, err := apiCfg.DB.GetBrowserCount90(r.Context(), params.DomainID)
		if err != nil {
			errChan <- err
			return
		}
		browserChan <- browser
	}()

	// Wait for all goroutines to finish
	var stats []database.GetNinetyDaysRow
	var os []database.GetOsCount90Row
	var device []database.GetDeviceCount90Row
	var browser []database.GetBrowserCount90Row

	for i := 0; i < 4; i++ {
		select {
		case result := <-statsChan:
			stats = result
		case result := <-osChan:
			os = result
		case result := <-deviceChan:
			device = result
		case result := <-browserChan:
			browser = result
		case err := <-errChan:
			respondWithError(w, http.StatusInternalServerError, "DB error: "+err.Error())
			return
		}
	}

	// Close channels
	close(statsChan)
	close(osChan)
	close(deviceChan)
	close(browserChan)
	close(errChan)

	
type response90 struct {
	Original []database.GetNinetyDaysRow       `json:"original"`
	Os       []database.GetOsCount90Row       `json:"os"`
	Device   []database.GetDeviceCount90Row   `json:"device"`
	Browser  []database.GetBrowserCount90Row  `json:"browser"`
}

	resp := response90{
		Original: stats,
		Os:       os,
		Device:   device,
		Browser:  browser,
	}

	elapsed := time.Since(startTime)
	fmt.Printf("my goroutine func took %s\n", elapsed)
	respondWithJson(w, 200, resp)
}


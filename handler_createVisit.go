package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/userAnalytics/internal/database"
)

// NOBODY LOOK HERE SECRET STUFF AND NOT SUPER BAD CODE
func extractDomain(referrer string) (string, error) {
	if referrer == "Direct visit" {
		return referrer, nil
	}
	parsedURL, err := url.Parse(referrer)
	if err != nil {
		return "", err
	}

	hostParts := strings.Split(parsedURL.Hostname(), ".")
	if len(hostParts) > 1 {
		return hostParts[len(hostParts)-2], nil
	}

	return parsedURL.Hostname(), nil
}
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
	cleanedRef, err := extractDomain(params.VisitFrom)
	if err != nil || len(cleanedRef) < 2 {
		cleanedRef = "Direct visit"
	}
	bounced := false
	if params.VisitDuration < 5 {
		bounced = true
	}
	if params.VisitDuration > 7500 {
		params.VisitDuration = 300
	}
	dbCtx := context.Background()
	// Asynchronously save the visit to the database
	go func() {
		_, err := apiCfg.DB.CreateVisit(dbCtx, database.CreateVisitParams{
			Createdat:     time.Now().UTC(),
			Visitorstatus: params.VisitStat,
			Visitduration: params.VisitDuration,
			Domain:        params.Domain,
			Visitfrom:     cleanedRef,
			Device:        params.Device,
			Os:            params.OS,
			Browser:       params.Browser,
			Bounce:        bounced,
		})
		if err != nil {
			fmt.Printf("error: %v \n", err)
		}
		var uniqueVisit int32 = 0
		if params.VisitStat == "new" {
			uniqueVisit = 1
		}
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

func getPageNameFromURL(rawURL string) (string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// Get the path from the URL
	path := parsedURL.Path
	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")
	// Split the path by slash to get individual components
	components := strings.Split(path, "/")
	// Reconstruct the path without UUIDs or other user-specific data
	var newPathComponents []string
	for _, component := range components {
		// Ignore if it looks like a UUID
		if len(component) == 36 && component[8] == '-' && component[13] == '-' && component[18] == '-' && component[23] == '-' {
			continue
		}
		// Add non-user-specific components to the new path
		newPathComponents = append(newPathComponents, component)
	}
	// Join the components to form the new path
	newPath := "/" + strings.Join(newPathComponents, "/")
	return newPath, nil
}

func (apiCfg *apiConfig) handlerCreatePageVisit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	type parameters struct {
		Domain uuid.UUID `json:"domain"`
		Page   string    `json:"page"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("error: %v \n", err)
		respondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	pageName, err := getPageNameFromURL(params.Page)
	if err != nil || len(pageName) == 0 {
		pageName = "/"
	}
	_, err = apiCfg.DB.CreatePageVisit(r.Context(), database.CreatePageVisitParams{
		Createdat: time.Now().UTC(),
		Domain:    params.Domain,
		Page:      pageName,
	})
	if err != nil {
		fmt.Printf("error with pagevisit: %v \n", err)
		return
	}
	respondWithJson(w, 200, "success")
}

// I CBA with this
// Testing speeds of different routines
func (apiCfg *apiConfig) handlerSevenVisits(w http.ResponseWriter, r *http.Request) {
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
	var wg sync.WaitGroup
	wg.Add(1)
	var stats []database.GetSevenDaysRow
	go func() {
		defer wg.Done()
		stats, err = apiCfg.DB.GetSevenDays(r.Context(), params.DomainID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
	}()
	var os []database.GetOsCount7Row
	wg.Add(1)
	go func() {
		defer wg.Done()
		os, err = apiCfg.DB.GetOsCount7(r.Context(), params.DomainID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
	}()
	var device []database.GetDeviceCount7Row
	wg.Add(1)
	go func() {
		defer wg.Done()
		device, err = apiCfg.DB.GetDeviceCount7(r.Context(), params.DomainID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
	}()
	var browser []database.GetBrowserCount7Row
	wg.Add(1)
	go func() {
		defer wg.Done()
		browser, err = apiCfg.DB.GetBrowserCount7(r.Context(), params.DomainID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
	}()
	var bounceRate float64
	wg.Add(1)
	go func() {
		defer wg.Done()
		bounceRate, err = apiCfg.DB.GetBounce7(r.Context(), params.DomainID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
	}()
	var pages []database.GetPages7Row
	wg.Add(1)
	go func() {
		defer wg.Done()
		pages, err = apiCfg.DB.GetPages7(r.Context(), params.DomainID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
	}()
	type response struct {
		Original   []database.GetSevenDaysRow     `json:"original"`
		Os         []database.GetOsCount7Row      `json:"os"`
		Device     []database.GetDeviceCount7Row  `json:"device"`
		Browser    []database.GetBrowserCount7Row `json:"browser"`
		BounceRate float64                        `json:"bounce"`
		Pages      []database.GetPages7Row        `json:"pages"`
	}
	wg.Wait()
	resp := response{
		Original:   stats,
		Os:         os,
		Device:     device,
		Browser:    browser,
		BounceRate: bounceRate,
		Pages:      pages,
	}
	elapsed := time.Since(startTime)
	fmt.Printf("my waitgroup funcs took %s\n", elapsed)
	respondWithJson(w, 200, resp)
}

// 30 days
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
	pages, err := apiCfg.DB.GetPages30(r.Context(), params.DomainID)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "DB error for pages")
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
	bounceRate, err := apiCfg.DB.GetBounce30(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}

	type response30 struct {
		Original   []database.GetLimitedCountRow   `json:"original"`
		Os         []database.GetOsCount30Row      `json:"os"`
		Device     []database.GetDeviceCount30Row  `json:"device"`
		Browser    []database.GetBrowserCount30Row `json:"browser"`
		BounceRate float64                         `json:"bounce"`
		Pages      []database.GetPages30Row        `json:"pages"`
	}

	resp := response30{
		Original:   stats,
		Os:         os,
		Device:     device,
		Browser:    browser,
		BounceRate: bounceRate,
		Pages:      pages,
	}
	elapsed := time.Since(startTime)
	fmt.Printf("my normal func took %s\n", elapsed)
	respondWithJson(w, 200, resp)
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
	os, err := apiCfg.DB.GetOsCount90(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error for Os")
		return
	}

	device, err := apiCfg.DB.GetDeviceCount90(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error for Device")
		return
	}

	browser, err := apiCfg.DB.GetBrowserCount90(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error for Browser")
		return
	}
	bounceRate, err := apiCfg.DB.GetBounce90(r.Context(), params.DomainID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}
	pages, err := apiCfg.DB.GetPages90(r.Context(), params.DomainID)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "DB error for pages")
		return
	}
	type response30 struct {
		Original   []database.GetNinetyDaysRow     `json:"original"`
		Os         []database.GetOsCount90Row      `json:"os"`
		Device     []database.GetDeviceCount90Row  `json:"device"`
		Browser    []database.GetBrowserCount90Row `json:"browser"`
		BounceRate float64                         `json:"bounce"`
		Pages      []database.GetPages90Row        `json:"pages"`
	}

	resp := response30{
		Original:   stats,
		Os:         os,
		Device:     device,
		Browser:    browser,
		BounceRate: bounceRate,
		Pages:      pages,
	}
	respondWithJson(w, 200, resp)

}

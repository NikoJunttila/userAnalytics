package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/nikojunttila/userAnalytics/internal/database"

  _"github.com/lib/pq"
)

type apiConfig struct{
  DB *database.Queries
}



func main() {
	r := chi.NewRouter()

	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port not found in env")
	}
	dbURL := os.Getenv("DB_URL")
	if portString == "" {
		log.Fatal("DB_URL is not found")
	}
  connection, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Fatal("cant connect to database",err)
  }

  apiCfg := apiConfig{
    DB: database.New(connection),
  }

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

  v1Router := chi.NewRouter()

  v1Router.Get("/healthz", handlerReadiness)
  v1Router.Get("/err",handlerErr)
  v1Router.Post("/users",apiCfg.handlerCreateUser)
  v1Router.Post("/visit",apiCfg.handlerCreateVisit)
  v1Router.Post("/visit/stats", apiCfg.handlerCountVisits)  
  v1Router.Post("/visit/seven", apiCfg.handlerSevenVisits)
  v1Router.Post("/visit/thirty", apiCfg.handlerLimitedVisits)
  v1Router.Post("/visit/ninety", apiCfg.handlerNinetyVisits)
  v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
  v1Router.Post("/domain",apiCfg.middlewareAuth(apiCfg.handlerGetDomain))
  v1Router.Post("/domains", apiCfg.middlewareAuth(apiCfg.handlerCreateDomain))
  v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerDomainFollowCreate))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerDomainFollowsGet))
  r.Mount("/v1", v1Router)

	fmt.Println("listening on port:" + portString)

	server := &http.Server{
		Handler: r,
		Addr:    ":" + portString,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

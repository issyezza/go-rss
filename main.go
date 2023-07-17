package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	// get env variables
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("no portString")
	}
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("no DB_URL")
	}

	// connect to db
	db, err := sql.Open("libsql", DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", DB_URL, err)
		os.Exit(1)
	}

	//  get data from db
	dbError := db.Ping()
	if dbError != nil {
		log.Fatal("no db ping")
	}

	{
		stmt, err := db.Prepare("SELECT * FROM users")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to prepare statement %s: %s", "SELECT * FROM users", err)
			os.Exit(1)
		}
		defer stmt.Close()
		rows, err := stmt.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute prepared statement %s: %s", "SELECT * FROM users", err)
			os.Exit(1)
		}
		for rows.Next() {
			var row struct {
				full_name  string
				email      string
				username   string
				id         int
				created_at int
			}
			if err := rows.Scan(&row.id, &row.username, &row.email, &row.full_name, &row.created_at); err != nil {
				fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
				os.Exit(1)
			}
			fmt.Println(row)
		}
		if err := rows.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "errors from query: %s", err)
			os.Exit(1)
		}
	}

	// setup server
	router := chi.NewRouter()
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/", handlerIndex)

	// create /v1/ready route
	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/error", handlerError)
	router.Mount("/v1", v1Router)

	// start server
	fmt.Printf("server running at: %v", portString)
	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatal("server is dead", err)
	}
}

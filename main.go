package main

import (
	companyHandler "backend/companyd/handler"
	companyRepo "backend/companyd/repository"
	"backend/companyd/usecase/company"
	userHandler "backend/userd/handler"
	"backend/userd/repository"
	"backend/userd/usecase/user"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	var err error
	var db *sql.DB
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "myuser")
	dbPassword := getEnv("DB_PASSWORD", "mypassword")
	dbName := getEnv("DB_NAME", "myapp")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Retry connection with exponential backoff
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Failed to open database: %v", err)
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to connect to database: %v (attempt %d/10)", err, i+1)
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatal("Could not connect to database after 10 attempts")
	}

	log.Println("Successfully connected to PostgreSQL!")

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Add logging middleware
	router.Use(loggingMiddleware)

	userdb := repository.NewRepository(db)
	// Register handlers with CORS middleware
	userHandler.RegisterHandlers(user.NewService(userdb), router)

	companydb := companyRepo.NewCompanyRepository(db)
	// Register handlers with CORS middleware
	companyHandler.RegisterHandlers(company.NewService(companydb), router)

	// Start server
	port := getEnv("PORT", "8080")
	serverAddr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

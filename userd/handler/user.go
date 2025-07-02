package userHandler

import (
	userPresenter "backend/userd/presenter"
	"backend/userd/usecase/user"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

// Helper function to get allowed origins from environment variable
func getAllowedOrigins() []string {
	// Default origins if environment variable is not set
	defaultOrigins := []string{}

	// Get from environment variable
	envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if envOrigins == "" {
		return defaultOrigins
	}

	// Split by comma and trim whitespace
	origins := strings.Split(envOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Get the origin from the request
		origin := r.Header.Get("Origin")

		// List of allowed origins from environment variable
		allowedOrigins := getAllowedOrigins()

		// Check if the origin is in the allowed list
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// If no origin matched, allow the request origin
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, ngrok-skip-browser-warning")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Content-Type", "application/json")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			log.Printf("Handling OPTIONS request for: %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

func UserLogin(service user.Usecase, w http.ResponseWriter, r *http.Request) {
	var loginRequest userPresenter.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	log.Printf("Login attempt for username: %s", loginRequest.Username)

	user, err := service.GetUserByUsername(loginRequest.Username, loginRequest.Password)
	if err != nil {
		log.Printf("Login error for username %s: %v", loginRequest.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid username or password",
		})
		return
	}

	// Mock user data - in real app, this would come from a database
	// users := map[string]struct {
	// 	ID        string `json:"id"`
	// 	Username  string `json:"username"`
	// 	Email     string `json:"email"`
	// 	Role      string `json:"role"`
	// 	CreatedAt string `json:"createdAt"`
	// }{
	// 	"admin": {
	// 		ID:        "1",
	// 		Username:  "admin",
	// 		Email:     "admin@company.com",
	// 		Role:      "Admin",
	// 		CreatedAt: "2024-01-01T00:00:00Z",
	// 	},
	// 	"manager": {
	// 		ID:        "2",
	// 		Username:  "manager",
	// 		Email:     "manager@company.com",
	// 		Role:      "Manager",
	// 		CreatedAt: "2024-01-01T00:00:00Z",
	// 	},
	// 	"officer": {
	// 		ID:        "3",
	// 		Username:  "officer",
	// 		Email:     "officer@company.com",
	// 		Role:      "Officer",
	// 		CreatedAt: "2024-01-01T00:00:00Z",
	// 	},
	// }

	// Check credentials
	// user, exists := users[loginRequest.Username]
	// if !exists || loginRequest.Password != "password" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(map[string]string{
	// 		"error": "Invalid username or password",
	// 	})
	// 	return
	// }

	// Return user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UserHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "user service is running",
		"status":  "running",
	}
	json.NewEncoder(w).Encode(response)
}

// Add a new endpoint to check database and list users for debugging
func UserDBTest(service user.Usecase, w http.ResponseWriter, r *http.Request) {
	users, err := service.ListUser()
	if err != nil {
		log.Printf("Database test error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"message":   "Database connection successful",
		"userCount": len(users),
		"users":     users,
	}
	json.NewEncoder(w).Encode(response)
}

func CreateUser(service user.Usecase, w http.ResponseWriter, r *http.Request) {
	var createRequest userPresenter.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	user, err := service.CreateUser(createRequest.Username, createRequest.Password, createRequest.Email, createRequest.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func ListUser(service user.Usecase, w http.ResponseWriter, r *http.Request) {
	users, err := service.ListUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(service user.Usecase, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]

	if !exists || id == "" {
		log.Printf("No ID provided in request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User ID is required",
		})
		return
	}

	log.Printf("Attempting to delete user with ID: %s", id)

	// Validate UUID format (case-insensitive)
	uuidRegex := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(id) {
		log.Printf("Invalid UUID format received: %s", id)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid UUID format",
		})
		return
	}

	err := service.DeleteUser(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Return user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}

func RegisterHandlers(service user.Usecase, router *mux.Router) {
	// Add CORS middleware to all routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enableCORS(next.ServeHTTP)(w, r)
		})
	})

	router.HandleFunc("/user/health", UserHealth).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/dbtest", func(w http.ResponseWriter, r *http.Request) {
		UserDBTest(service, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		UserLogin(service, w, r)
	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		CreateUser(service, w, r)
	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/list", func(w http.ResponseWriter, r *http.Request) {
		ListUser(service, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteUser(service, w, r)
	}).Methods("DELETE", "OPTIONS")
}

package companyHandler

import (
	companyPresenter "backend/companyd/presenter"
	"backend/companyd/usecase/company"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Helper function to get allowed origins from environment variable
func getAllowedOrigins() []string {
	// Default origins if environment variable is not set
	defaultOrigins := []string{
		"https://0f22-2402-3a80-1325-cd70-dd05-94a2-213-dd84.ngrok-free.app",
		"https://place-pro-platform-88.vercel.app",
		"https://localhost:8081",
		"http://localhost:8081",
	}

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

func CompanyHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "company service is running",
		"status":  "running",
	}
	json.NewEncoder(w).Encode(response)
}

func CreateCompany(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	var createRequest companyPresenter.CreateCompany
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	company, err := service.CreateCompany(
		createRequest.CompanyName,
		createRequest.CompanyAddress,
		createRequest.Drive,
		createRequest.TypeOfDrive,
		createRequest.FollowUp,
		strconv.FormatBool(createRequest.IsContacted),
		createRequest.Remarks,
		createRequest.ContactDetails,
		createRequest.Hr1Details,
		createRequest.Hr2Details,
		createRequest.Package,
		createRequest.AssignedOfficer,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

func ListCompanies(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	companies, err := service.ListCompanies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}

func ListCompaniesByUsername(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	username, exists := vars["id"]

	if !exists || username == "" {
		log.Printf("No username provided in request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Company username is required",
		})
		return
	}

	companies, err := service.ListCompaniesByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}

func DeleteCompany(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, exists := vars["id"]

	if !exists || id == "" {
		log.Printf("No ID provided in request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Company ID is required",
		})
		return
	}

	log.Printf("Attempting to delete company with ID: %s", id)

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

	err := service.DeleteCompany(id)
	if err != nil {
		log.Printf("Error deleting company: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Company deleted successfully",
	})
}

func UpdateCompany(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, exists := vars["id"]

	if !exists || id == "" {
		log.Printf("No ID provided in request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Company ID is required",
		})
		return
	}

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

	var updateRequest companyPresenter.CreateCompany
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	company, err := service.UpdateCompany(
		id,
		updateRequest.CompanyName,
		updateRequest.CompanyAddress,
		updateRequest.Drive,
		updateRequest.TypeOfDrive,
		updateRequest.FollowUp,
		strconv.FormatBool(updateRequest.IsContacted),
		updateRequest.Remarks,
		updateRequest.ContactDetails,
		updateRequest.Hr1Details,
		updateRequest.Hr2Details,
		updateRequest.Package,
		updateRequest.AssignedOfficer,
	)
	if err != nil {
		log.Printf("Error updating company: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

func CreateCompanyTemp(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createRequest companyPresenter.CreateCompanyTemp
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	companyTemp, err := service.CreateCompanyTemp(
		createRequest.CompanyID,
		createRequest.CompanyName,
		createRequest.CompanyAddress,
		createRequest.Drive,
		createRequest.TypeOfDrive,
		createRequest.FollowUp,
		strconv.FormatBool(createRequest.IsContacted),
		createRequest.Remarks,
		createRequest.ContactDetails,
		createRequest.Hr1Details,
		createRequest.Hr2Details,
		createRequest.Package,
		createRequest.AssignedOfficer,
		createRequest.CreatedBy,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companyTemp)
}

func ListCompanyTemps(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	companyTemps, err := service.ListCompanyTemps()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companyTemps)
}

func UpdateCompanyTempStatus(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Company temp ID is required",
		})
		return
	}

	var updateRequest struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	err := service.UpdateCompanyTempStatus(id, updateRequest.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Company temp status updated successfully",
	})
}

func ApproveCompanyTemp(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Company temp ID is required",
		})
		return
	}

	err := service.ApproveCompanyTemp(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Company update approved and applied successfully",
	})
}

func CreateEvent(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createRequest struct {
		Date        string `json:"date"`
		Type        string `json:"type"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatedBy   string `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if createRequest.Date == "" || createRequest.Type == "" || createRequest.Title == "" || createRequest.CreatedBy == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing required fields: date, type, title, and created_by are required",
		})
		return
	}

	event, err := service.CreateEvent(
		createRequest.Date,
		createRequest.Type,
		createRequest.Title,
		createRequest.Description,
		createRequest.CreatedBy,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

func ListEvents(service company.Usecase, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	events, err := service.ListEvents()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Format events for frontend
	formattedEvents := make([]map[string]interface{}, len(events))
	for i, event := range events {
		formattedEvents[i] = map[string]interface{}{
			"id":          event.ID,
			"date":        event.Date,
			"type":        event.Type,
			"title":       event.Title,
			"description": event.Description,
			"createdBy":   event.CreatedBy,
			"createdAt":   event.CreatedAt,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formattedEvents)
}

func RegisterHandlers(service company.Usecase, router *mux.Router) {
	// Add CORS middleware to all routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enableCORS(next.ServeHTTP)(w, r)
		})
	})

	router.HandleFunc("/company/health", CompanyHealth).Methods("GET", "OPTIONS")
	router.HandleFunc("/company/create", func(w http.ResponseWriter, r *http.Request) {
		CreateCompany(service, w, r)
	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/company/list", func(w http.ResponseWriter, r *http.Request) {
		ListCompanies(service, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/company/list/{id}", func(w http.ResponseWriter, r *http.Request) {
		ListCompaniesByUsername(service, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/company/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteCompany(service, w, r)
	}).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/company/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateCompany(service, w, r)
	}).Methods("PUT", "OPTIONS")
	router.HandleFunc("/company/temp/update", func(w http.ResponseWriter, r *http.Request) {
		CreateCompanyTemp(service, w, r)
	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/company/temp/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		CreateCompanyTemp(service, w, r)
	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/company/temp/list", func(w http.ResponseWriter, r *http.Request) {
		ListCompanyTemps(service, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/company/temp/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateCompanyTempStatus(service, w, r)
	}).Methods("PUT", "OPTIONS")
	router.HandleFunc("/company/temp/approve/{id}", func(w http.ResponseWriter, r *http.Request) {
		ApproveCompanyTemp(service, w, r)
	}).Methods("PUT", "OPTIONS")
	router.HandleFunc("/event/create", func(w http.ResponseWriter, r *http.Request) {
		CreateEvent(service, w, r)
	}).Methods("POST", "OPTIONS")
	router.HandleFunc("/event/list", func(w http.ResponseWriter, r *http.Request) {
		ListEvents(service, w, r)
	}).Methods("GET", "OPTIONS")
}

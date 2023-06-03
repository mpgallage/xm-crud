package main

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var db *gorm.DB

type CompanyType string

const (
	Corporations       CompanyType = "Corporations"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "SoleProprietorship"
)

type Company struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name          string      `gorm:"type:varchar(15);not null;unique"`
	Description   string      `gorm:"type:varchar(3000)"`
	EmployeeCount int         `gorm:"type:int;not null"`
	Registered    bool        `gorm:"type:boolean;not null"`
	Type          CompanyType `gorm:"type:varchar(50);not null"` //Corporations | NonProfit | Cooperative | Sole Proprietorship
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_ARGS"))
	if err != nil {
		log.Fatal("Error connecting to database.", err)
		panic("Failed to connect database")
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			log.Error("Error closing database!", err)
		}
	}(db)

	tx := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if tx.Error != nil {
		log.Fatal("Error creating extension.", tx.Error)
		return
	}

	db.AutoMigrate(&Company{})

	r := mux.NewRouter()
	r.HandleFunc("/company", createCompanyHandler).Methods("POST")
	r.HandleFunc("/company/{id}", getCompanyHandler).Methods("GET")
	r.HandleFunc("/company/{id}", updateCompanyHandler).Methods("PATCH")
	r.HandleFunc("/company/{id}", deleteCompanyHandler).Methods("DELETE")

	http.Handle("/", r)
	err = http.ListenAndServe(":8080", jsonContentTypeMiddleware(r))
	if err != nil {
		log.Fatal("Error starting server.", err)
		return
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func createCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var company Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidCompanyType(company.Type) {
		http.Error(w, "Invalid company type", http.StatusBadRequest)
		return
	}

	if err = db.Create(&company).Error; err != nil {
		http.Error(w, "Invalid values for properties.", http.StatusBadRequest)
		return
	}
	// TODO: produce event to Kafka
	err = json.NewEncoder(w).Encode(company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func getCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var company Company
	if err := db.First(&company, "id = ?", id).Error; err != nil {
		http.Error(w, "No record found.", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func updateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var updatedCompany Company
	err = json.NewDecoder(r.Body).Decode(&updatedCompany)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidCompanyType(updatedCompany.Type) {
		http.Error(w, "Invalid company type", http.StatusBadRequest)
		return
	}

	var company Company
	if err := db.First(&company, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = db.Model(&company).Updates(updatedCompany).Error; err != nil {
		http.Error(w, "Invalid values for properties.", http.StatusBadRequest)
		return
	}

	// TODO: produce event to Kafka
	err = json.NewEncoder(w).Encode(company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func deleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var company Company
	if err := db.First(&company, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	db.Delete(&company)
	// TODO: produce event to Kafka
	w.WriteHeader(http.StatusNoContent)
}

var companyTypes = map[CompanyType]bool{
	Corporations:       true,
	NonProfit:          true,
	Cooperative:        true,
	SoleProprietorship: true,
}

func isValidCompanyType(t CompanyType) bool {
	return companyTypes[t]
}

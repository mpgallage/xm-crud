package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/mpgallage/xmcrud/database"
	"github.com/mpgallage/xmcrud/events"
	"github.com/mpgallage/xmcrud/models"
	"math/rand"
	"net/http"
)

func CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var company models.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !models.IsValidCompanyType(company.Type) {
		http.Error(w, "Invalid company type", http.StatusBadRequest)
		return
	}

	if err = database.Get().Create(&company).Error; err != nil {
		http.Error(w, "Invalid values for properties.", http.StatusBadRequest)
		return
	}

	events.ProduceKafka(fmt.Sprintf("create-%s-%v", company.ID.String(), rand.Int()), company)
	err = json.NewEncoder(w).Encode(company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var company models.Company
	if err := database.Get().First(&company, "id = ?", id).Error; err != nil {
		http.Error(w, "No record found.", http.StatusNotFound)
		return
	}
	events.ProduceKafka(fmt.Sprintf("get-%s-%v", company.ID.String(), rand.Int()), company)
	_ = json.NewEncoder(w).Encode(company)
}

func UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var updatedCompany models.Company
	err = json.NewDecoder(r.Body).Decode(&updatedCompany)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !models.IsValidCompanyType(updatedCompany.Type) {
		http.Error(w, "Invalid company type", http.StatusBadRequest)
		return
	}

	var company models.Company
	if err := database.Get().First(&company, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err = database.Get().Model(&company).Updates(updatedCompany).Error; err != nil {
		http.Error(w, "Invalid values for properties.", http.StatusBadRequest)
		return
	}

	events.ProduceKafka(fmt.Sprintf("update-%s-%v", updatedCompany.ID.String(), rand.Int()), company)
	_ = json.NewEncoder(w).Encode(company)
}

func DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var company models.Company
	if err := database.Get().First(&company, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	database.Get().Delete(&company)

	events.ProduceKafka(fmt.Sprintf("delete-%s-%v", company.ID.String(), rand.Int()), company)
	w.WriteHeader(http.StatusNoContent)
}

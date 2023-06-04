package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mpgallage/xmcrud/models"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http"
	"testing"
)

var serverUrl = "http://localhost:8081/company"

func getToken() string {
	authReq :=
		"{" +
			"\"Username\": \"malaka\"," +
			"\"Password\": \"password\"" +
			"}"
	response, _ := http.Post(authUrl, "application/json", bytes.NewReader([]byte(authReq)))
	body, _ := io.ReadAll(response.Body)

	var authResp models.JwtToken
	_ = json.Unmarshal(body, &authResp)
	return authResp.Token
}

func TestCreateCompany(t *testing.T) {
	token := getToken()
	name := fmt.Sprintf("ABC-%v", rand.Intn(100000))
	company, status := createCompany(name, models.Corporations, token)

	assert.Equal(t, status, 200)
	assert.NotEmpty(t, company.ID)
	assert.Equal(t, company.Name, name)
}

func TestCreateInvalidCompany(t *testing.T) {
	token := getToken()
	name := fmt.Sprintf("ABC-%v", rand.Intn(100000))
	_, status := createCompany(name, "Invalid Type", token)

	assert.Equal(t, status, 400)
}

func TestGetCompany(t *testing.T) {
	token := getToken()
	name := fmt.Sprintf("ABC-%v", rand.Intn(100000))
	company, _ := createCompany(name, models.Corporations, token)

	newCompany, status := getCompany(company.ID.String(), token)
	assert.Equal(t, status, 200)
	assert.Equal(t, company.ID, newCompany.ID)
	assert.Equal(t, company.Name, newCompany.Name)
}

func TestGetCompanyNonExistentID(t *testing.T) {
	token := getToken()
	_, status := getCompany("fe105b10-2139-46f5-a08b-a0b69594a595", token)
	assert.Equal(t, status, 404)
}

func TestGetCompanyInvalidID(t *testing.T) {
	token := getToken()
	_, status := getCompany("fe105b10-2139-46f5-a08b", token)
	assert.Equal(t, status, 400)
}

func TestUpdateCompany(t *testing.T) {
	token := getToken()
	name := fmt.Sprintf("ABC-%v", rand.Intn(100000))
	company, _ := createCompany(name, models.Corporations, token)

	newName := fmt.Sprintf("ABC-%v", rand.Intn(100000))
	jsonReq := []byte(fmt.Sprintf("{"+
		"\"Name\": \"%s\","+
		"\"Description\": \"This is a description of company.\","+
		"\"EmployeeCount\": 150,"+
		"\"Registered\": true,"+
		"\"Type\": \"%s\""+
		"}", newName, models.SoleProprietorship))

	newCompany, status := updateCompany(company.ID.String(), jsonReq, token)
	assert.Equal(t, status, 200)
	assert.Equal(t, company.ID, newCompany.ID)
	assert.Equal(t, newCompany.Name, newName)
	assert.Equal(t, newCompany.Type, models.SoleProprietorship)
}

func TestDeleteCompany(t *testing.T) {
	token := getToken()
	name := fmt.Sprintf("ABC-%v", rand.Intn(100000))
	company, _ := createCompany(name, models.Corporations, token)

	newCompany, status := getCompany(company.ID.String(), token)
	assert.Equal(t, status, 200)
	assert.Equal(t, company.ID, newCompany.ID)
	assert.Equal(t, company.Name, newCompany.Name)

	status = deleteCompany(company.ID.String(), token)
	assert.Equal(t, status, 204)

	_, status = getCompany(company.ID.String(), token)
	assert.Equal(t, status, 404)
}

func TestDeleteCompanyNonExistentID(t *testing.T) {
	token := getToken()
	status := deleteCompany("fe105b10-2139-46f5-a08b-a0b69594a595", token)
	assert.Equal(t, status, 404)
}

func createCompany(name string, compType models.CompanyType, token string) (models.Company, int) {
	jsonReq := []byte(fmt.Sprintf("{"+
		"\"Name\": \"%s\","+
		"\"Description\": \"This is a description of company.\","+
		"\"EmployeeCount\": 150,"+
		"\"Registered\": true,"+
		"\"Type\": \"%s\""+
		"}", name, compType))
	req, _ := http.NewRequest("POST", serverUrl, bytes.NewBuffer(jsonReq))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var company models.Company
	err := json.Unmarshal(body, &company)
	if err != nil {
		return models.Company{}, resp.StatusCode
	}
	return company, resp.StatusCode
}

func getCompany(id string, token string) (models.Company, int) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", serverUrl, id), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var company models.Company
	err := json.Unmarshal(body, &company)
	if err != nil {
		return models.Company{}, resp.StatusCode
	}
	return company, resp.StatusCode
}

func updateCompany(id string, jsonReq []byte, token string) (models.Company, int) {
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", serverUrl, id), bytes.NewBuffer(jsonReq))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var company models.Company
	err := json.Unmarshal(body, &company)
	if err != nil {
		return models.Company{}, resp.StatusCode
	}
	return company, resp.StatusCode
}

func deleteCompany(id string, token string) int {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", serverUrl, id), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, _ := client.Do(req)

	return resp.StatusCode
}

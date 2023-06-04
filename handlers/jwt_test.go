package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/mpgallage/xmcrud/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

var authUrl = "http://localhost:8081/authenticate"

func TestValidJWT(t *testing.T) {
	authReq :=
		"{" +
			"\"Username\": \"malaka\"," +
			"\"Password\": \"password\"" +
			"}"
	response, _ := http.Post(authUrl, "application/json", bytes.NewReader([]byte(authReq)))
	body, err := io.ReadAll(response.Body)

	assert.NoError(t, err)
	var authResp models.JwtToken
	_ = json.Unmarshal(body, &authResp)
	assert.NotEmpty(t, authResp.Token)
}

func TestInvalidJWT(t *testing.T) {
	authReq := ""
	response, _ := http.Post(authUrl, "application/json", bytes.NewReader([]byte(authReq)))
	body, err := io.ReadAll(response.Body)

	assert.NoError(t, err)
	var authResp models.Exception
	_ = json.Unmarshal(body, &authResp)
	assert.Equal(t, authResp.Message, "Invalid input.")
}

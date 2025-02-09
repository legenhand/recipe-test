package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/legenhand/recipe-test/src/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSubmitEmail(t *testing.T) {
	config.Cfg = &config.Config{
		JWTSecret: "test-secret-key",
	}
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/submit-email", SubmitEmail)

	requestBody := map[string]string{
		"email": "firmansyah720.fs@gmail.com",
	}

	jsonValue, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/auth/submit-email", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}

	assert.Contains(t, response, "message")
	assert.Equal(t, "Magic link sent. Please check your email.", response["message"])
}

func TestMagicLink(t *testing.T) {
	config.Cfg = &config.Config{
		JWTSecret: "test-secret-key",
	}
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/auth/magic-link", MagicLink)

	validToken, _ := GenerateMagicLinkToken("test@gmail.com)", 15*time.Minute)
	req, _ := http.NewRequest("GET", "/auth/magic-link?token="+validToken, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Authentication successful", response["message"])
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
}

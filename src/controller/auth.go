package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/legenhand/recipe-test/src/config"
	"github.com/legenhand/recipe-test/src/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SubmitEmail(c *gin.Context) {
	var payload model.SubmitEmailRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email provided"})
		return
	}

	tokenString, err := GenerateMagicLinkToken(payload.Email, 15*time.Minute)
	if err != nil {
		log.Printf("Failed to generate magic link token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate magic link token"})
		return
	}

	magicLinkURL := fmt.Sprintf("%s/auth/magic-link?token=%s", config.Cfg.BaseUrl, tokenString)

	log.Printf("Send magic link to %s: %s", payload.Email, magicLinkURL)

	c.JSON(http.StatusOK, gin.H{"message": "Magic link sent. Please check your email."})
}

func MagicLink(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token parameter"})
		return
	}

	email, err := ValidateMagicLinkToken(tokenString)
	if err != nil {
		log.Printf("Failed to validate magic link token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	authClaims := model.MagicLinkClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // auth token valid for 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	authTokenString, err := authToken.SignedString([]byte(config.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"token":   authTokenString,
	})
}

func GenerateMagicLinkToken(email string, duration time.Duration) (string, error) {
	expiry := time.Now().Add(duration).Unix()

	data := fmt.Sprintf("%s:%d", email, expiry)

	mac := hmac.New(sha256.New, []byte(config.Cfg.JWTSecret))
	_, err := mac.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("failed to write data to HMAC: %w", err)
	}
	signature := hex.EncodeToString(mac.Sum(nil))

	token := fmt.Sprintf("%s:%d:%s", email, expiry, signature)

	encodedToken := base64.URLEncoding.EncodeToString([]byte(token))

	return encodedToken, nil
}

func ValidateMagicLinkToken(tokenString string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to decode token: %w", err)
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}

	email := parts[0]
	expiryStr := parts[1]
	signatureProvided := parts[2]

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid expiry timestamp: %w", err)
	}

	if time.Now().Unix() > expiry {
		return "", errors.New("token expired")
	}

	data := fmt.Sprintf("%s:%s", email, expiryStr)
	mac := hmac.New(sha256.New, []byte(config.Cfg.JWTSecret))
	_, err = mac.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("failed to write data to HMAC: %w", err)
	}
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signatureProvided), []byte(expectedSignature)) {
		return "", errors.New("invalid token signature")
	}

	return email, nil
}

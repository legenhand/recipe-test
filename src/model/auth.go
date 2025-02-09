package model

import "github.com/golang-jwt/jwt/v5"

type SubmitEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SubmitEmailResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type MagicLinkResponse struct {
	Email string               `json:"email"`
	Token jwt.RegisteredClaims `json:"token"`
}

type MagicLinkClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

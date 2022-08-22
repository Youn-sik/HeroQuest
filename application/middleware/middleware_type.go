package middleware

import "github.com/dgrijalva/jwt-go/v4"

type AuthTokenClaims struct {
	ID                 string `json:"id"` // 유저 ID
	jwt.StandardClaims        // 표준 토큰 Claims
}

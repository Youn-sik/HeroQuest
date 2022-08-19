package user

import "github.com/dgrijalva/jwt-go/v4"

type Send_data struct {
	result bool
	errStr string
}

type User struct {
	Id           string `json:"id,omitempty"`
	Account      string `json:"account,omitempty"`
	Password     string `json:"password,omitempty"`
	Name         string `json:"name,omitempty"`
	TokenBalance int    `json:"token_balance,omitempty"`
	QId          string `json:"qid,omitempty"`
}

type CreateUserReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type ModifyUserReq struct {
	Account      string `json:"account"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	TokenBalance int    `json:"token_balance"`
	QId          string `json:"qid"`
}

type DeleteUserReq struct {
	Id string `json:"id"`
}

type InfoUserReq struct {
	Id string `json:"id"`
}

type AuthTokenClaims struct {
	ID                 string `json:"id"` // 유저 ID
	jwt.StandardClaims        // 표준 토큰 Claims
}

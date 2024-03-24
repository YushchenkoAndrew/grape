package types

import "github.com/golang-jwt/jwt/v5"

type AccessClaim struct {
	jwt.RegisteredClaims
	UID    string `json:"uid"`
	RID    string `json:"rid"`
	UserId int64  `json:"user_id"`
	Exp    int64  `json:"exp"`
}

type RefreshClaim struct {
	jwt.RegisteredClaims
	UID string `json:"uid"`
	Exp int64  `json:"exp"`
}

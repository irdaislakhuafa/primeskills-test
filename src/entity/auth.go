package entity

import "github.com/golang-jwt/jwt/v5"

type (
	AuthJWTClaims struct {
		UID string `json:"uid"`
		jwt.RegisteredClaims
	}
)

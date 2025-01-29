package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type (
	TokenServiceParams struct {
		JwtSecret string
	}

	TokenService struct {
		jwtSecret []byte
	}

	Claims struct {
		Email string `json:"email"`
		jwt.StandardClaims
	}
)

func NewTokenService(params TokenServiceParams) *TokenService {
	return &TokenService{
		jwtSecret: []byte(params.JwtSecret),
	}
}

func (t *TokenService) GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "gymratz-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.jwtSecret)
}

func (t *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return t.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

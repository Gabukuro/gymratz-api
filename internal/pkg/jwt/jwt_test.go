package jwt_test

import (
	"testing"
	"time"

	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	goJwt "github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

var testJwtKey = "test-jwt-key"
var testEmail string = "test@example.com"

func TestTokenService(t *testing.T) {
	t.Parallel()

	tokenService := jwt.NewTokenService(jwt.TokenServiceParams{
		JwtSecret: testJwtKey,
	})

	t.Run("should generate and validate a token to ensure it's correctly generated", func(t *testing.T) {
		token, err := tokenService.GenerateToken(testEmail)
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.IsType(t, "", token)
		assert.NotEmpty(t, token)

		claims, err := tokenService.ValidateToken(token)
		assert.Nil(t, err)
		assert.Equal(t, testEmail, claims.Email)
	})

	t.Run("should return an error when the token is invalid", func(t *testing.T) {
		invalidToken := "invalid token"

		claims, err := tokenService.ValidateToken(invalidToken)
		assert.NotNil(t, err)
		assert.Nil(t, claims)
	})

	t.Run("should return an error when the token is expired", func(t *testing.T) {
		token, err := tokenService.GenerateToken(testEmail)
		assert.Nil(t, err)
		assert.NotNil(t, token)
		assert.IsType(t, "", token)
		assert.NotEmpty(t, token)

		expiredToken := goJwt.NewWithClaims(goJwt.SigningMethodHS256, &jwt.Claims{
			Email: testEmail,
			StandardClaims: goJwt.StandardClaims{
				ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			},
		})
		expiredTokenString, err := expiredToken.SignedString([]byte(testJwtKey))
		assert.Nil(t, err)
		assert.NotNil(t, expiredTokenString)

		_, err = tokenService.ValidateToken(expiredTokenString)
		assert.NotNil(t, err)
	})
}

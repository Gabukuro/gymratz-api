package testhelper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	internalJWT "github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/golang-jwt/jwt"
)

func RunRequest(
	setup *setup.Setup,
	method string,
	path string,
	body any,
	header map[string]string,
) (*http.Response, error) {
	req := httptest.NewRequest(
		method,
		path,
		parseBodyToStringReader(body),
	)

	if _, ok := header["Authorization"]; !ok {
		req.Header.Add("Authorization", generateAuthToken(setup.EnvVariables.JWTSecret))
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Add(key, value)
	}

	return setup.App.Test(req, -1)
}

func parseBodyToStringReader(requestBody any) *strings.Reader {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(string(jsonBody))
}

func generateAuthToken(secret string) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := internalJWT.Claims{
		Email: "test@email.com",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "gymratz-api-test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return "Bearer " + signedToken
}

func ParseSuccessResponseBody[data any](body io.ReadCloser) response.SuccessResponse[data] {
	var responseBody response.SuccessResponse[data]

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		panic(err)
	}

	return responseBody
}

func ParsePaginationResponseBody[data any](body io.ReadCloser) response.PaginationResponse[data] {
	var responseBody response.PaginationResponse[data]

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		panic(err)
	}

	return responseBody
}

func ParseErrorResponseBody(body io.ReadCloser) response.ErrorResponse {
	var responseBody response.ErrorResponse

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		panic(err)
	}

	return responseBody
}

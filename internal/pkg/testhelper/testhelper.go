package testhelper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
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

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TestResponse struct {
	StatusCode   int
	ResponseBody map[string]string
}

func DoPost(requestBody interface{}, endpoint string, contentType string) TestResponse {
	requestBodyReader := StructToReader(requestBody)

	resp, err := http.Post(endpoint, contentType, requestBodyReader)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error making post: %s", err))
	}

	responseBody := ReadCloserToMap(resp.Body)
	return TestResponse{resp.StatusCode, responseBody}
}

func StructToReader(requestBody interface{}) io.Reader {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error marshalling json data: %s", err))
	}

	requestBodyReader := bytes.NewReader(jsonBody)
	return requestBodyReader
}

func ReadCloserToMap(rc io.ReadCloser) map[string]string {
	responseBodyBytes, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading response: %s", err))
	}

	var responseBody map[string]string
	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error unmarshalling json data: %s", err))
	}

	return responseBody
}

func DoPostWithCookie(requestBody interface{}, endpoint string, contentType string, cookieName string, cookieValue string) TestResponse {
	client := &http.Client{}
	requestBodyReader := StructToReader(requestBody)
	req, err := http.NewRequest("POST", endpoint, requestBodyReader)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating request: %s", err))
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", cookieName, cookieValue))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error executing request: %s", err))
	}

	responseBody := ReadCloserToMap(resp.Body)
	return TestResponse{resp.StatusCode, responseBody}
}

func DoPostWithCookieNoBodyInResponse(requestBody interface{}, endpoint string, contentType string, cookieName string, cookieValue string) TestResponse {
	client := &http.Client{}
	requestBodyReader := StructToReader(requestBody)
	req, err := http.NewRequest("POST", endpoint, requestBodyReader)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error creating request: %s", err))
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", cookieName, cookieValue))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error executing request: %s", err))
	}

	return TestResponse{resp.StatusCode, nil}
}

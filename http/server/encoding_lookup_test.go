package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func IssueEncodePasswordRequest(application *ApplicationServer, password string) (persistenceResult, error) {
	formData := url.Values{"password": {password}}
	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))
	request.Form = formData
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(application.EncodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	responseContent := responseRecorder.Body.Bytes()
	var result persistenceResult

	if err := json.Unmarshal(responseContent, &result); err != nil {
		return result, err
	}

	return result, nil
}

func IssueEncodedPasswordLookupRequest(application *ApplicationServer, url string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(application.LookupEncodingByID)

	handler.ServeHTTP(responseRecorder, request)
	return responseRecorder
}

func TestLookupEncodingByIDWithValidIDReturnsHTTPOKStatus(t *testing.T) {
	application := NewAppServer()
	result, err := IssueEncodePasswordRequest(application, "P@ssW0rd")

	if err != nil {
		t.Error(err)
	}

	response := IssueEncodedPasswordLookupRequest(application, result.URL)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Got %v, Wanted %v", status, http.StatusOK)
	}
}

func TestLookupEncodingByIDWithValidIDReturnsPlainTextContentType(t *testing.T) {
	application := NewAppServer()
	result, err := IssueEncodePasswordRequest(application, "P@ssW0rd")

	if err != nil {
		t.Error(err)
	}

	response := IssueEncodedPasswordLookupRequest(application, result.URL)

	if contentType := response.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
		t.Errorf("Handler returned wrong Content Type. Got %v, Wanted %v", contentType, "text/plain; charset=utf-8")
	}
}

func TestLookupEncodingByIDWithValidIDReturnsEncodedPassword(t *testing.T) {
	application := NewAppServer()
	result, err := IssueEncodePasswordRequest(application, "P@ssW0rd!")

	if err != nil {
		t.Error(err)
	}

	time.Sleep(6 * time.Second)

	response := IssueEncodedPasswordLookupRequest(application, result.URL)
	encodedPassword := response.Body.String()

	if encodedPassword != "62+j0x1/W8bCgSgF3YggMtf+AfOqb28xuOXvKvTXBs8iDZDwQci9cGBiNdHvHHyywclJeKIhPWoftStSNJdf5g==" {
		t.Errorf("Handler returned wrong encoded Password. Got %v, Wanted %v",
			encodedPassword,
			"62+j0x1/W8bCgSgF3YggMtf+AfOqb28xuOXvKvTXBs8iDZDwQci9cGBiNdHvHHyywclJeKIhPWoftStSNJdf5g==")
	}
}

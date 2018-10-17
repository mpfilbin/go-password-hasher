package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestEncodeAndPersistWithValidFormDataReturnsStatusAccepted(t *testing.T) {
	application := NewAppServer()
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(application.encodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code. Got %v, Wanted %v", status, http.StatusAccepted)
	}
}

func TestEncodeAndPersistWithValidFormDataReturnsJSONContentType(t *testing.T) {
	application := NewAppServer()
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(application.encodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != ContentTypeJson {
		t.Errorf("Handler returned wrong content type. Got %v, Wanted %v", contentType, ContentTypeJson)
	}
}

func TestEncodeAndPersistWithValidFormDataReturnsPersistenceResultWithLookupURL(t *testing.T) {
	application := NewAppServer()
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(application.encodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	responseContent := responseRecorder.Body.Bytes()
	var result encodingResult

	if err := json.Unmarshal(responseContent, &result); err != nil {
		t.Error(err)
	}

	if result.URL != "/hash/1" {
		t.Errorf("Handler returned wrong payload. Expected %v to match %v", result.URL, "/hash/1")
	}
}

func TestEncodeAndPersistWithValidFormDataReturnsPersistenceResultWithTimeAvailable(t *testing.T) {
	application := NewAppServer()
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(application.encodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	responseContent := responseRecorder.Body.Bytes()
	var result encodingResult

	if err := json.Unmarshal(responseContent, &result); err != nil {
		t.Error(err)
	}

	if result.TimeAvailable == "" {
		t.Errorf("Handler return wrong payload. TimeAvailable should not be empty")
	}
}

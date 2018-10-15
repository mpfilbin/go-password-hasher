package password

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)


func TestEncodeAndPersistWithValidFormDataReturnsStatusAccepted(t *testing.T) {
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(EncodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code. Got %v, Wanted %v", status, http.StatusAccepted)
	}
}

func TestEncodeAndPersistWithValidFormDataReturnsJSONContentType(t *testing.T) {
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(EncodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	if contentType := responseRecorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong content type. Got %v, Wanted %v", contentType, "application/json")
	}
}

func TestEncodeAndPersistWithValidFormDataReturnsPersistenceResultWithLookupURL(t *testing.T) {
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(EncodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	responseContent := responseRecorder.Body.Bytes()
	var result persistenceResult

	if err := json.Unmarshal(responseContent, &result); err != nil {
		t.Error(err)
	}

	ok, err := regexp.Match("\\/hash\\/\\d+", []byte(result.URL))

	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Errorf("Handler returned wrong payload. Expected %v to match /hash/{number}", result.URL)
	}
}

func TestEncodeAndPersistWithValidFormDataReturnsPersistenceResultWithTimeAvailable(t *testing.T) {
	formData := url.Values{}
	formData.Add("password", "P@ssW0rd")

	request := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(formData.Encode()))

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(EncodeAndPersist)

	handler.ServeHTTP(responseRecorder, request)

	responseContent := responseRecorder.Body.Bytes()
	var result persistenceResult

	if err := json.Unmarshal(responseContent, &result); err != nil {
		t.Error(err)
	}

	if result.TimeAvailable == "" {
		t.Errorf("Handler return wrong payload. TimeAvailable should not be empty")
	}
}
package services

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleEmoji(t *testing.T) {
	requestBody := bytes.NewBuffer([]byte(`{"emoji": "😀"}`))
	req, err := http.NewRequest("POST", "/emoji", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleEmoji)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `Received emoji: 😀`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

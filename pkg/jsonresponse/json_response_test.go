package jsonresponse

import (
	_ "bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusOK(t *testing.T) {
	// Создаем тестовый HTTP-респондер
	w := httptest.NewRecorder()
	response := MessageResp{Message: "Test message"}

	// Вызываем функцию StatusOK
	StatusOK(w, response)

	// Проверяем статус-код
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	// Проверяем тело ответа
	var responseBody MessageResp
	err := json.NewDecoder(w.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if responseBody.Message != response.Message {
		t.Errorf("expected message %s, got %s", response.Message, responseBody.Message)
	}
}

func TestError(t *testing.T) {
	// Создаем тестовый HTTP-респондер
	w := httptest.NewRecorder()
	errMessage := "some error occurred"

	// Вызываем функцию Error
	Error(w, errors.New(errMessage), http.StatusBadRequest)

	// Проверяем статус-код
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, res.StatusCode)
	}

	// Проверяем заголовки
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
	}

	// Проверяем тело ответа
	var errorResponse map[string]string
	err := json.NewDecoder(w.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if errorResponse["errors"] != errMessage {
		t.Errorf("expected error message %s, got %s", errMessage, errorResponse["errors"])
	}
}

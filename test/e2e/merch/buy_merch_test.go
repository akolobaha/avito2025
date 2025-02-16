package merch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func TestBuyMerch(t *testing.T) {
	// Данные для аутентификации
	authData := AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	// Шаг 1: Получение JWT токена
	token, err := getAuthToken(authData)
	if err != nil {
		t.Fatalf("Ошибка аутентификации: %v", err)
	}

	// Шаг 2: Попытка покупки предмета
	item := "t-shirt" // Подставьте идентификатор товара
	err = buyItem(item, token)
	if err != nil {
		t.Fatalf("Ошибка покупки предмета: %v", err)
	}

	t.Logf("Успешно куплен предмет: %s", item)
}

// Функция для получения JWT токена
func getAuthToken(authData AuthRequest) (string, error) {
	authURL := "http://0.0.0.0:8080/api/auth"
	jsonData, err := json.Marshal(authData)
	if err != nil {
		return "", fmt.Errorf("ошибка маршалинга данных: %v", err)
	}

	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неудачная аутентификация, статус код: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения тела ответа: %v", err)
	}

	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return "", fmt.Errorf("ошибка разбора ответа: %v", err)
	}

	return authResp.Token, nil
}

// Функция для покупки предмета
func buyItem(item, token string) error {
	buyURL := fmt.Sprintf("http://0.0.0.0:8080/api/buy/%s", item)

	req, err := http.NewRequest("GET", buyURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Добавляем JWT токен в заголовки
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось купить предмет, статус код: %d", resp.StatusCode)
	}

	return nil
}

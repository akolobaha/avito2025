package coins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// Структуры для аутентификации и ответа
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func TestSendCoins(t *testing.T) {
	// Данные для аутентификации
	authDataFrom := AuthRequest{
		Username: "testuser1",
		Password: "testpassword",
	}

	authDataTo := AuthRequest{
		Username: "testuser2",
		Password: "testpassword",
	}

	// Шаг 1: Получение JWT токена
	token, err := getAuthToken(authDataFrom)
	if err != nil {
		t.Fatalf("Ошибка аутентификации отправителя: %v", err)
	}

	// Шаг 2: Создадим пользователя, которму отправим деньги
	_, err = getAuthToken(authDataTo)
	if err != nil {
		t.Fatalf("Ошибка аутентификации получателя: %v", err)
	}

	// Шаг 3: Перевод монет другому пользователю
	sendData := SendCoinRequest{
		ToUser: authDataTo.Username,
		Amount: 1,
	}
	err = sendCoins(sendData, token)
	if err != nil {
		t.Fatalf("Ошибка перевода монет: %v", err)
	}

	t.Logf("Успешно переведено %d монет пользователю %s", sendData.Amount, sendData.ToUser)
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

// Функция для перевода монет
func sendCoins(sendData SendCoinRequest, token string) error {
	sendCoinURL := "http://0.0.0.0:8080/api/sendCoin"

	// Создание запроса
	jsonData, err := json.Marshal(sendData)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга данных: %v", err)
	}

	req, err := http.NewRequest("POST", sendCoinURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Добавляем JWT токен в заголовки
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось перевести монеты, статус код: %d", resp.StatusCode)
	}

	return nil
}

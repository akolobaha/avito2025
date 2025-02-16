package info

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

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []Transaction `json:"received"`
	Sent     []Transaction `json:"sent"`
}

type Transaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

func TestGetInfo(t *testing.T) {
	// Данные для аутентификации
	authData := AuthRequest{
		Username: "testuser3",
		Password: "testpassword",
	}

	// Шаг 1: Получение JWT токена
	token, err := getAuthToken(authData)
	if err != nil {
		t.Fatalf("Ошибка аутентификации: %v", err)
	}

	// Шаг 2: Получение информации о монетах, инвентаре и истории транзакций
	info, err := getInfo(token)
	if err != nil {
		t.Fatalf("Ошибка получения информации: %v", err)
	}

	// Проверка значений в ответе
	t.Logf("Информация о монетах: %d", info.Coins)
	t.Logf("Инвентарь: %+v", info.Inventory)
	t.Logf("История монет - Полученные: %+v, Отправленные: %+v", info.CoinHistory.Received, info.CoinHistory.Sent)

	// Пример проверки значения монет
	if info.Coins < 0 {
		t.Errorf("Неверное количество монет: %d", info.Coins)
	}
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

// Функция для получения информации о монетах, инвентаре и истории транзакций
func getInfo(token string) (InfoResponse, error) {
	infoURL := "http://0.0.0.0:8080/api/info"

	req, err := http.NewRequest("GET", infoURL, nil)
	if err != nil {
		return InfoResponse{}, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Добавляем JWT токен в заголовки
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return InfoResponse{}, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return InfoResponse{}, fmt.Errorf("не удалось получить информацию, статус код: %d", resp.StatusCode)
	}

	var infoResp InfoResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return InfoResponse{}, fmt.Errorf("ошибка чтения тела ответа: %v", err)
	}

	err = json.Unmarshal(body, &infoResp)
	if err != nil {
		return InfoResponse{}, fmt.Errorf("ошибка разбора ответа: %v", err)
	}

	return infoResp, nil
}

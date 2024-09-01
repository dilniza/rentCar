package pkg

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"rent-car/config"
	"sync"
	"time"
)

func NullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func Duration(fD, tD string) (duration float64, err error) {
	fromDate, err := time.Parse(time.RFC3339, fD)
	if err != nil {
		return 0, err
	}

	toDate, err := time.Parse(time.RFC3339, tD)
	if err != nil {
		return 0, err
	}

	duration = float64(toDate.Sub(fromDate).Hours() / 24)
	return duration, nil
}

func GetSerialId(n *int) string {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	*n++
	result := fmt.Sprintf("Or-%08d", *n)
	return result
}

func TelegramBotFunc(msg interface{}) (string, error) {

	botToken := config.BotToken
	chatID := config.ChatID
	
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	message := string(messageBytes)

	payload := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: chatID,
		Text:   message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	success := "Response Status:" + resp.Status
	return success, nil
}

func GenerateOTP() int {

	return rand.Intn(900000) + 100000
}

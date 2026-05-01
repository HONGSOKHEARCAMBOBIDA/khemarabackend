package helper

import (
	"bytes"
	"encoding/json"
	"net/http"

	"strconv"
)

type TelegramMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func SendTelegramMessage(message string) error {
	botToken := "8619235719:AAHGC0_NVzIm9Gchc0tzD7nbGZ9Grv5xyG4"
	chatID := "-1003747714579"

	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	body := TelegramMessage{
		ChatID:    parseInt64(chatID),
		Text:      message,
		ParseMode: "HTML",
	}

	jsonData, _ := json.Marshal(body)

	_, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	return err
}

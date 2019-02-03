package telegram

import (
	"net/http"
	"strconv"
)

type TlgBot struct {
	token  string
	userId uint64
}

const apiTelegram = "https://api.telegram.org/bot"

func NewTelegram(token string, userId uint64) *TlgBot {
	return &TlgBot{
		token:  token,
		userId: userId,
	}
}

func (bot *TlgBot) SendMessage(message string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", apiTelegram+bot.token+"/sendMessage", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.FormatUint(bot.userId, 10))
	q.Add("text", "[To-Do Family]: "+message)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

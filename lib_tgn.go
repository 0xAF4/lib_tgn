package lib_tgn

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TelegramNotifier struct {
	token  string
	admins []string
}

const (
	TelegramSendURL = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=HTML"
)

func New(token string, adms []string) (*TelegramNotifier, error) {
	if len(adms) == 0 {
		return nil, errors.New("Укажите админов бота")
	}

	return &TelegramNotifier{
		token:  token,
		admins: adms,
	}, nil
}

func (b *TelegramNotifier) Notify(message string) error {
	for _, admin := range b.admins {
		url := fmt.Sprintf(TelegramSendURL, b.token, admin, message)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if !strings.Contains(string(body), "message") {
			return errors.New("NOT SUCCESS")
		}
	}
	return nil
}

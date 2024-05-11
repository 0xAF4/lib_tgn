package lib_tgn

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type TelegramNotifier struct {
	token  string
	prefix string
	admins *[]string
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

const (
	sendMessageURL = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=HTML"
	getUpdates     = "https://api.telegram.org/bot%s/getUpdates"
	regGetChats    = `"chat":(.*?),`
)

func New(token string, pref string, adms *[]string) (*TelegramNotifier, error) {
	if len(*adms) == 0 {
		return nil, errors.New("Укажите админов бота")
	}

	body, err := SendHttpGet(fmt.Sprintf(getUpdates, token))
	if err != nil {
		return nil, err
	}

	chats := &[]Chat{}
	re := regexp.MustCompile(`"chat":(.*?),`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		if len(match) > 1 {
			data := new(Chat)
			if err := json.Unmarshal([]byte(match[1]), data); err == nil {
				*chats = append(*chats, *data)
			}
		}
	}

	for i, admin := range *adms {
		if !hasLetters(admin) {
			continue
		}
		if id := FindChatIDbyUsername(chats, admin); id != "" {
			(*adms)[i] = id
		}
	}

	return &TelegramNotifier{
		token:  token,
		prefix: pref,
		admins: adms,
	}, nil
}

func (b *TelegramNotifier) Notify(message string) error {
	for _, admin := range *b.admins {
		if hasLetters(admin) {
			continue
		}
		if _, err := SendHttpGet(fmt.Sprintf(sendMessageURL, b.token, admin, b.prefix+"\n"+message)); err != nil {
			return err
		}
	}
	return nil
}

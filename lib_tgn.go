package lib_tgn

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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

const (
	LevelInfo    = "🔵INFO🔵"
	LevelSuccess = "🟢SUCCESS🟢"
	LevelWarning = "🟡WARNING🟡"
	LevelError   = "🔴ERROR🔴"
)

// Возвращает чаты в виде JSON объектов
// chats, _ := l_tgn.GetChats(conf.Telegram.Token)
// fmt.Println(chats)

func GetChats(token string) ([]Chat, error) {
	body, err := SendHttpGet(fmt.Sprintf(getUpdates, token))
	if err != nil {
		return nil, err
	}

	chats := &[]Chat{}
	re := regexp.MustCompile(`"chat":{(.*?)},`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		if len(match) > 1 {
			data := new(Chat)
			if err := json.Unmarshal([]byte("{"+match[1]+"}"), data); err == nil {
				*chats = append(*chats, *data)
			}
		}
	}

	return *chats, nil
}

func New(token string, pref string, adms *[]string) (*TelegramNotifier, error) {
	if len(*adms) == 0 {
		return nil, errors.New("Укажите админов бота")
	}

	var iadms []string
	for _, admin := range *adms {
		if hasLetters(admin) {
			continue
		}
		iadms = append(iadms, admin)
	}

	return &TelegramNotifier{
		token:  token,
		prefix: pref,
		admins: &iadms,
	}, nil
}

func (b *TelegramNotifier) Notify(message string) error {
	for _, admin := range *b.admins {
		if hasLetters(admin) {
			continue
		}
		if _, err := SendHttpGet(fmt.Sprintf(sendMessageURL, b.token, admin, url.QueryEscape(b.prefix+"\n"+message))); err != nil {
			return err
		}
	}
	return nil
}

func (b *TelegramNotifier) NotifyWithLevel(message string, level string) error {
	for _, admin := range *b.admins {
		if hasLetters(admin) {
			continue
		}
		if _, err := SendHttpGet(fmt.Sprintf(sendMessageURL, b.token, admin, url.QueryEscape(b.prefix+"\n"+level+"\n"+message))); err != nil {
			return err
		}
	}
	return nil
}

func (b *TelegramNotifier) AsyncNotify(message string) {
	go b.Notify(message)
}

func (b *TelegramNotifier) AsyncNotifyWithLevel(message string, level string) {
	go b.NotifyWithLevel(message, level)
}

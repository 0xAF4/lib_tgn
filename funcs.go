package lib_tgn

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func hasLetters(str string) bool {
	for _, char := range str {
		if unicode.IsLetter(char) {
			return true
		}
	}
	return false
}

func FindChatIDbyUsername(arr *[]Chat, username string) string {
	for _, chat := range *arr {
		if chat.Username == username {
			return strconv.Itoa(chat.ID)
		}
	}
	return ""
}

func SendHttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if !strings.Contains(string(body), "message") {
		return "", errors.New("NOT SUCCESS")
	}

	return string(body), nil
}

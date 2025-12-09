package webapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// VerifyInitData checks initData HMAC signature and returns TelegramUser when valid.
func VerifyInitData(initData string, botToken string) (*TelegramUser, error) {
	if initData == "" {
		return nil, fmt.Errorf("empty initData")
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("parse initData: %w", err)
	}

	hash := values.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash is missing")
	}
	values.Del("hash")

	var dataPairs []string
	for key, vals := range values {
		value := vals[0]
		dataPairs = append(dataPairs, fmt.Sprintf("%s=%s", key, value))
	}

	sort.Strings(dataPairs)

	dataCheckString := strings.Join(dataPairs, "\n")

	secretKey := hmacSHA256([]byte("WebAppData"), []byte(botToken))
	calculatedHashBytes := hmacSHA256(secretKey, []byte(dataCheckString))
	calculatedHash := hex.EncodeToString(calculatedHashBytes)

	if !hmac.Equal([]byte(calculatedHash), []byte(hash)) {
		return nil, fmt.Errorf("invalid hash")
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, fmt.Errorf("user field is missing")
	}

	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}

	return &user, nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

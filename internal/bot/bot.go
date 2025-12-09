package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// BotClient provides minimal Telegram Bot API usage: sendMessage + long polling.
type BotClient struct {
	token  string
	apiURL string
	client *http.Client
	offset int
}

func NewBotClient(token string) *BotClient {
	return &BotClient{
		token:  token,
		apiURL: "https://api.telegram.org",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (b *BotClient) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("%s/bot%s/sendMessage", b.apiURL, b.token)

	payload := sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal sendMessage payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create sendMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("sendMessage http error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sendMessage non-200: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// --- Long Polling ---

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type getUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func (b *BotClient) StartLongPolling() error {
	for {
		updates, err := b.getUpdates()
		if err != nil {
			log.Printf("getUpdates error: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, u := range updates {
			if u.UpdateID >= int64(b.offset) {
				b.offset = int(u.UpdateID) + 1
			}
			if u.Message == nil {
				continue
			}
			b.handleMessage(u.Message)
		}
	}
}

func (b *BotClient) getUpdates() ([]Update, error) {
	url := fmt.Sprintf("%s/bot%s/getUpdates?timeout=20&offset=%d", b.apiURL, b.token, b.offset)

	resp, err := b.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getUpdates http error: %w", err)
	}
	defer resp.Body.Close()

	var result getUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode getUpdates response: %w", err)
	}

	if !result.OK {
		return nil, fmt.Errorf("getUpdates not ok")
	}

	return result.Result, nil
}

func (b *BotClient) handleMessage(m *Message) {
	if m.Text == "/start" {
		_ = b.SendMessage(m.Chat.ID, "Привет! 👋\n\nНажми кнопку Mini App внизу, чтобы сыграть в крестики-нолики 💕")
		return
	}
}

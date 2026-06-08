package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// TelegramSender mengirim pesan teks ke user via Telegram Bot API.
// "to" adalah Chat ID Telegram (disimpan di kolom telegram_chat_id tabel users).
type TelegramSender struct {
	BotToken string
}

// NewTelegramSender membuat instance TelegramSender dari env TELEGRAM_BOT_TOKEN.
// Mengembalikan error jika token tidak dikonfigurasi.
func NewTelegramSender() (*TelegramSender, error) {
	token := GetEnv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN tidak dikonfigurasi")
	}
	return &TelegramSender{BotToken: token}, nil
}

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// Send mengirim pesan ke Chat ID tertentu.
// Mendukung Markdown (bold **text**, italic _text_, link [teks](url)).
// Mengembalikan response string dari Telegram API dan error jika gagal.
func (t *TelegramSender) Send(chatID, text string) (string, error) {
	if chatID == "" {
		return "", fmt.Errorf("telegram chat_id kosong")
	}
	if text == "" {
		return "", fmt.Errorf("pesan telegram kosong")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	// Ensure we don't forward raw HTML tags (e.g. <p>) to Telegram.
	// Telegram is configured to parse Markdown; sending raw HTML tags
	// results in visible tags in the message bubble. Strip tags here.
	cleanText := stripHTMLTags(text)

	payload := telegramPayload{
		ChatID:    chatID,
		Text:      cleanText,
		ParseMode: "Markdown",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("telegram marshal payload: %w", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("telegram HTTP send error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	respText := string(respBody)

	if resp.StatusCode >= 400 {
		log.Printf("[TELEGRAM] Error %d: %s", resp.StatusCode, respText)
		return respText, fmt.Errorf("telegram API error %d: %s", resp.StatusCode, respText)
	}

	log.Printf("[TELEGRAM] Sent to chatID=%s status=%d", chatID, resp.StatusCode)
	return fmt.Sprintf("telegram:ok:%d", resp.StatusCode), nil
}

// stripHTMLTags removes simple HTML tags like <p>, <br>, <strong>, etc.
// It is intentionally simple (state machine) to avoid importing regex.
func stripHTMLTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			if inTag {
				inTag = false
			} else {
				b.WriteRune(r)
			}
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}
	// Collapse multiple consecutive newlines and trim spaces
	out := strings.ReplaceAll(b.String(), "\r", "")
	out = strings.TrimSpace(out)
	return out
}

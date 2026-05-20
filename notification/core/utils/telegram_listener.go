package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// StartTelegramListener menjalankan goroutine yang mendengarkan pesan dari pengguna
// via long-polling (getUpdates). Ini digunakan agar bot bisa merespon perintah /start
// dengan mengirimkan Chat ID kepada pengguna.
func StartTelegramListener(ctx context.Context) {
	token := GetEnv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("[TELEGRAM LISTENER] Bot token tidak dikonfigurasi, listener dinonaktifkan.")
		return
	}

	go func() {
		offset := 0
		client := &http.Client{Timeout: 35 * time.Second} // Timeout lebih besar dari long-polling timeout
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", token)

		log.Println("[TELEGRAM LISTENER] Berjalan di background (long-polling)...")

		for {
			select {
			case <-ctx.Done():
				log.Println("[TELEGRAM LISTENER] Berhenti.")
				return
			default:
				// Lanjut polling
			}

			// Polling dengan timeout 30 detik
			url := fmt.Sprintf("%s?offset=%d&timeout=30", apiURL, offset)
			resp, err := client.Get(url)
			if err != nil {
				// Tunggu sebentar jika terjadi error koneksi
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}

			var result struct {
				Ok     bool `json:"ok"`
				Result []struct {
					UpdateID int `json:"update_id"`
					Message  struct {
						MessageID int `json:"message_id"`
						Text      string `json:"text"`
						Chat      struct {
							ID   int64  `json:"id"`
							Type string `json:"type"`
						} `json:"chat"`
						From struct {
							FirstName string `json:"first_name"`
						} `json:"from"`
					} `json:"message"`
				} `json:"result"`
			}

			if err := json.Unmarshal(body, &result); err != nil || !result.Ok {
				time.Sleep(2 * time.Second)
				continue
			}

			for _, update := range result.Result {
				offset = update.UpdateID + 1

				if update.Message.Text != "" && update.Message.Chat.Type == "private" {
					chatID := update.Message.Chat.ID
					text := update.Message.Text
					firstName := update.Message.From.FirstName

					// Cek jika perintah adalah /start
					if text == "/start" {
						replyText := fmt.Sprintf("Halo %s! 👋\n\nChat ID Anda adalah:\n`%d`\n\nSilakan salin angka di atas dan masukkan ke halaman Pengaturan di aplikasi ThinkTala untuk menghubungkan notifikasi Telegram.", firstName, chatID)
						
						// Kirim balasan menggunakan utils.TelegramSender yang sudah ada
						tg, _ := NewTelegramSender()
						if tg != nil {
							tg.Send(fmt.Sprintf("%d", chatID), replyText)
						}
					}
				}
			}
		}
	}()
}

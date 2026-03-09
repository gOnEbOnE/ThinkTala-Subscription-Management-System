package system

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

func Env(req string) string {
	// 1. Cek OS environment dulu (Railway inject vars ke OS)
	if val := os.Getenv(req); val != "" {
		return val
	}

	// 2. Fallback ke .env file jika ada
	var f *os.File
	var err error

	f, err = os.Open(".env")

	if err != nil {
		// .env tidak ada, return empty (sudah cek OS env di atas)
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key := strings.Split(scanner.Text(), "=")
		if strings.TrimSpace(key[0]) == req {
			return strings.TrimSpace(key[1])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	return ""
}

var (
	refererOnce sync.Once
	refererList map[string]bool
)

// LoadReferers membaca referer whitelist dari file referers.json
func LoadReferers() map[string]bool {
	refererOnce.Do(func() {
		data, err := os.ReadFile("referers.json")
		if err != nil {
			refererList = make(map[string]bool)
			return
		}
		_ = json.Unmarshal(data, &refererList)
	})
	return refererList
}

/*
	Copyright © 2024 - 2025. PT Arunika Tala Archipelago
	Developed by Muhammad Abror
	For more info, please visit https://arunikatala.co.id
*/

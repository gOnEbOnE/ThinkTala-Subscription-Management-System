package utils

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

var (
	envCache = make(map[string]string)
	envMu    sync.Mutex
	envLoaded bool
)

// LoadEnv reads .env file and caches it.
// If the file is not found, callers can try another path.
func LoadEnv(filepath string) {
	envMu.Lock()
	defer envMu.Unlock()

	if envLoaded {
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envCache[key] = value
		}
	}

	envLoaded = true
}

// GetEnv retrieves value from cache or os environment
func GetEnv(key string, fallback ...string) string {
	// 1. Cek OS Environment dulu (untuk Docker/Kubernetes)
	if val := os.Getenv(key); val != "" {
		return val
	}
	// 2. Cek Cache .env file
	if val, ok := envCache[key]; ok {
		return val
	}
	// 3. Fallback
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

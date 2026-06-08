package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// ============================================================
// 1. KONFIGURASI GLOBAL & PEMETAAN PAKET (TIERING)
// ============================================================
const (
	RedisHost     = "72.61.215.108:6379"
	RedisPass     = "M3l4t1cn"
	RedisProdHash = "market_signal_production"
)

var tierPackages = map[string][]string{
	// Basic: 1 major pair saja sebagai teaser
	"basic": {"EURUSD"},

	// Pro: Major Forex pairs
	"pro": {
		"EURUSD", "GBPUSD", "AUDUSD", "USDCAD", "USDCHF",
		"USDJPY", "NZDUSD", "EURGBP", "EURJPY", "GBPJPY",
		"AUDJPY", "CADJPY", "CHFJPY", "EURAUD", "EURCAD",
		"EURCHF", "GBPAUD", "GBPCAD", "GBPCHF", "AUDCAD",
		"AUDNZD", "AUDCHF", "NZDJPY", "NZDCAD", "NZDCHF",
		"CADCHF",
	},

	// Enterprise: Forex + Crypto + Metals + Commodities
	"enterprise": {
		// Forex Major
		"EURUSD", "GBPUSD", "AUDUSD", "USDCAD", "USDCHF",
		"USDJPY", "NZDUSD", "EURGBP", "EURJPY", "GBPJPY",
		"AUDJPY", "CADJPY", "CHFJPY", "EURAUD", "EURCAD",
		"EURCHF", "GBPAUD", "GBPCAD", "GBPCHF", "AUDCAD",
		"AUDNZD", "AUDCHF", "NZDJPY", "NZDCAD", "NZDCHF",
		"CADCHF", "EURNZD", "GBPNZD",
		// Forex Exotic
		"USDTRY", "USDZAR", "USDMXN", "USDNOK", "USDSEK",
		"USDPLN", "EURPLN", "EURSEK", "EURNOK",
		// Crypto Top
		"BTCUSD", "ETHUSD", "SOLUSD", "BNBUSD", "XRPUSD",
		"ADAUSD", "DOTUSD", "AVXUSD", "LTCUSD", "XLMUSD",
		"LINKUSD", "UNIUSD", "AAVEUSD", "MATICUSD",
		// Metals & Commodities
		"XAUUSD", "XAGUSD", "XPTUSD", "XPDUSD",
		"XTIUSD", "XBRUSD", "XNGUSD",
	},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	signalCache = make(map[string]map[string]interface{})
	cacheMutex  = &sync.RWMutex{}
)

// ============================================================
// 2. CLIENT PROFILE STRUCT
// ============================================================
type Client struct {
	conn *websocket.Conn
	tier string
}

func isPairAllowed(tier, pair string) bool {
	allowedPairs, exists := tierPackages[tier]
	if !exists {
		return false
	}
	for _, p := range allowedPairs {
		if p == pair {
			return true
		}
	}
	return false
}

// ============================================================
// 3. WEBSOCKET HUB (DIFILTER PER TIER)
// ============================================================
type WSHub struct {
	clients   map[*Client]bool
	broadcast chan map[string]map[string]interface{}
	mu        sync.Mutex
}

func newWSHub() *WSHub {
	return &WSHub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan map[string]map[string]interface{}, 256),
	}
}

func (h *WSHub) run() {
	for currentCache := range h.broadcast {
		h.mu.Lock()

		tierPayloads := make(map[string][]byte)
		for tier := range tierPackages {
			var filteredSignals []interface{}
			for pairName, data := range currentCache {
				if isPairAllowed(tier, pairName) {
					filteredSignals = append(filteredSignals, data)
				}
			}

			wsMsg := map[string]interface{}{
				"type":    "active_signals",
				"payload": filteredSignals,
				"count":   len(filteredSignals),
			}
			if bytes, err := json.Marshal(wsMsg); err == nil {
				tierPayloads[tier] = bytes
			}
		}

		for client := range h.clients {
			msgBytes, exists := tierPayloads[client.tier]
			if !exists || len(msgBytes) == 0 {
				continue
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				log.Println("❌ WS Write error:", err)
				client.conn.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

func (h *WSHub) add(c *Client)    { h.mu.Lock(); h.clients[c] = true; h.mu.Unlock() }
func (h *WSHub) remove(c *Client) { h.mu.Lock(); delete(h.clients, c); h.mu.Unlock() }

// ============================================================
// 5. REDIS CONSUMER (MENELAN SEMUA DATA KE MEMORI)
// ============================================================
func listenToProductionHash(hub *WSHub) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: RedisHost, Password: RedisPass, DB: 0})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("❌ Gagal connect Redis: %v\n", err)
		return
	}
	log.Printf("✅ Menelan SEMUA data dari Redis Hash: [%s]\n", RedisProdHash)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	localState := make(map[string]string)
	var lastPricesStr string

	for range ticker.C {
		// 1. Read signals from Hash
		states, err := rdb.HGetAll(ctx, RedisProdHash).Result()
		if err != nil {
			log.Println("⚠️ Error membaca Redis Hash:", err)
			continue
		}

		hasUpdates := false

		// 2. Read latest prices from Stream
		var pricesMap map[string]map[string]interface{}
		pricesUpdated := false
		priceMsgs, err := rdb.XRevRangeN(ctx, "price_realtime", "+", "-", 1).Result()
		if err == nil && len(priceMsgs) > 0 {
			if pricesJSON, ok := priceMsgs[0].Values["prices"].(string); ok {
				if pricesJSON != lastPricesStr {
					lastPricesStr = pricesJSON
					pricesUpdated = true
					json.Unmarshal([]byte(pricesJSON), &pricesMap)
				}
			}
		}

		// Update signals cache
		for pairName, jsonStr := range states {
			var signalData map[string]interface{}
			cacheMutex.RLock()
			cachedData, exists := signalCache[pairName]
			cacheMutex.RUnlock()

			isNewState := localState[pairName] != jsonStr
			if isNewState {
				localState[pairName] = jsonStr
				hasUpdates = true
				if err := json.Unmarshal([]byte(jsonStr), &signalData); err != nil {
					continue
				}
				signalData["pair"] = pairName
				if rawDataStr, ok := signalData["data"].(string); ok && rawDataStr != "" {
					var parsedRawData map[string]interface{}
					if errParse := json.Unmarshal([]byte(rawDataStr), &parsedRawData); errParse == nil {
						signalData["data"] = parsedRawData
					}
				}
			} else if exists {
				// Copy existing cached data to avoid modifying during read
				signalData = make(map[string]interface{})
				cacheMutex.RLock()
				for k, v := range cachedData {
					signalData[k] = v
				}
				cacheMutex.RUnlock()
			} else {
				continue
			}

			// Inject real-time prices if available
			if pricesMap != nil {
				if pairPrices, ok := pricesMap[pairName]; ok {
					if bid, ok := pairPrices["bid"]; ok {
						signalData["bid_price"] = fmt.Sprintf("%v", bid)
						signalData["bid"] = bid
					}
					if ask, ok := pairPrices["ask"]; ok {
						signalData["ask_price"] = fmt.Sprintf("%v", ask)
						signalData["ask"] = ask
					}
				}
			}

			cacheMutex.Lock()
			signalCache[pairName] = signalData
			cacheMutex.Unlock()
		}

		if hasUpdates || (pricesUpdated && len(signalCache) > 0) {
			cacheMutex.RLock()
			cacheSnapshot := make(map[string]map[string]interface{})
			for k, v := range signalCache {
				cacheSnapshot[k] = v
			}
			cacheMutex.RUnlock()

			// Broadcast snapshot
			select {
			case hub.broadcast <- cacheSnapshot:
			default:
			}
		}
	}
}

// wsHandler di-bind ke /ws di gateway
func handleWebSocket(hub *WSHub, w http.ResponseWriter, r *http.Request) {
	clientTier := r.URL.Query().Get("tier")
	clientTier = strings.ToLower(clientTier)
	
	if _, valid := tierPackages[clientTier]; !valid {
		clientTier = "basic"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade WS Error:", err)
		return
	}

	clientObj := &Client{conn: conn, tier: clientTier}
	hub.add(clientObj)
	log.Printf("📡 Klien Terhubung: %s | Paket: [%s]\n", r.RemoteAddr, strings.ToUpper(clientTier))

	cacheMutex.RLock()
	var initSignals []interface{}
	for pairName, data := range signalCache {
		if isPairAllowed(clientTier, pairName) {
			initSignals = append(initSignals, data)
		}
	}
	cacheMutex.RUnlock()

	if len(initSignals) > 0 {
		initPayload := map[string]interface{}{
			"type":    "active_signals",
			"payload": initSignals,
			"count":   len(initSignals),
		}
		if initBytes, err := json.Marshal(initPayload); err == nil {
			conn.WriteMessage(websocket.TextMessage, initBytes)
		}
	}

	defer func() {
		hub.remove(clientObj)
		conn.Close()
		log.Printf("🔌 Klien Terputus: %s\n", r.RemoteAddr)
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// InitWSEngine memulai WSHub dan Redis listener
func InitWSEngine() *WSHub {
	hub := newWSHub()
	go hub.run()
	go listenToProductionHash(hub)
	return hub
}

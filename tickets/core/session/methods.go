package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tickets/core/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions" // Pastikan import ini ada untuk &sessions.Options{}
)

// Helper Internal untuk memastikan SessionManager tidak NIL
func getManager() (*Manager, error) {
	if SessionManager == nil {
		return nil, errors.New("SessionManager is NIL. Did you call session.Init() in main.go?")
	}
	return SessionManager, nil
}

// Set Session
func Set(w http.ResponseWriter, r *http.Request, key string, value any) (string, error) {
	// 1. Cek Manager
	mgr, err := getManager()
	if err != nil {
		fmt.Println("[PANIC AVOIDED] " + err.Error())
		return "", err
	}

	// A. STATEFUL (Cookie/File/RedisStore)
	if mgr.Config.Driver == "cookie" || mgr.Config.Driver == "file" || mgr.Config.Driver == "redis" {
		// 2. Cek Store
		if mgr.Store == nil {
			return "", errors.New("Session Store is NIL. Check SESSION_DRIVER in .env")
		}

		// 3. Get Session dengan Error Handling
		sess, err := mgr.Store.Get(r, mgr.Config.CookieName)

		// Gorilla sessions kadang mereturn error jika cookie corrupt/invalid,
		// tapi tetap mengembalikan objek session baru.
		// Kita hanya panic jika objek session benar-benar NIL.
		if sess == nil {
			msg := "Session object is NIL from Store.Get"
			if err != nil {
				msg += ": " + err.Error()
			}
			return "", errors.New(msg)
		}

		// 4. FIX PANIC: Pastikan Options tidak NIL sebelum diakses
		if sess.Options == nil {
			sess.Options = &sessions.Options{}
		}

		// 5. Set Value & Options
		sess.Values[key] = value

		// Re-apply Options
		sess.Options.MaxAge = mgr.Config.SessionLife
		sess.Options.Path = mgr.Config.Path
		sess.Options.Domain = mgr.Config.Domain
		sess.Options.Secure = mgr.Config.Secure
		sess.Options.HttpOnly = mgr.Config.HttpOnly
		sess.Options.SameSite = mgr.Config.SameSite

		errSave := sess.Save(r, w)

		fmt.Println("ada kok", sess.Values[key])
		return "", errSave
	}

	// B. STATELESS (JWT Pure)
	if mgr.Config.Driver == "stateless" {
		claims := jwt.MapClaims{
			"exp":  time.Now().Add(time.Duration(mgr.Config.SessionLife) * time.Second).Unix(),
			"data": value,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString([]byte(mgr.Config.SecretKey))
	}

	// C. STATELESS WITH REDIS (Hybrid)
	if mgr.Config.Driver == "stateless_with_redis" {
		sessionID := uuid.NewString()

		// Serialize Data
		jsonVal, _ := json.Marshal(value)

		// Gunakan utils.RedisSet agar aman dan reuse koneksi
		err := utils.RedisSet(
			context.Background(),
			"sess:"+sessionID,
			jsonVal,
			time.Duration(mgr.Config.SessionLife)*time.Second,
		)

		if err != nil {
			return "", err
		}

		// Buat Token Referensi
		claims := jwt.MapClaims{
			"exp": time.Now().Add(time.Duration(mgr.Config.SessionLife) * time.Second).Unix(),
			"sid": sessionID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString([]byte(mgr.Config.SecretKey))
	}

	return "", errors.New("invalid session driver")
}

// Get Session
func Get(r *http.Request, key string) any {
	mgr, err := getManager()
	if err != nil {
		return nil
	}

	// A. STATEFUL
	if mgr.Config.Driver == "cookie" || mgr.Config.Driver == "file" || mgr.Config.Driver == "redis" {
		if mgr.Store == nil {
			return nil
		}

		sess, _ := mgr.Store.Get(r, mgr.Config.CookieName)
		// Cek nil sess agar tidak panic
		if sess == nil {
			return nil
		}

		return sess.Values[key]
	}

	// Helper ambil token
	tokenString := extractToken(r)
	if tokenString == "" {
		return nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(mgr.Config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	// B. STATELESS
	if mgr.Config.Driver == "stateless" {
		data, ok := claims["data"].(map[string]any)
		if !ok {
			return nil
		}
		return data[key]
	}

	// C. STATELESS WITH REDIS
	if mgr.Config.Driver == "stateless_with_redis" {
		sid, ok := claims["sid"].(string)
		if !ok {
			return nil
		}

		valStr, err := utils.RedisGet(context.Background(), "sess:"+sid)
		if err != nil {
			return nil
		}

		var data map[string]any
		json.Unmarshal([]byte(valStr), &data)
		return data[key]
	}

	return nil
}

// Destroy / Logout
func Destroy(w http.ResponseWriter, r *http.Request) error {
	mgr, err := getManager()
	if err != nil {
		return err
	}

	// A. Stateful
	if mgr.Config.Driver == "cookie" || mgr.Config.Driver == "file" || mgr.Config.Driver == "redis" {
		if mgr.Store == nil {
			return nil
		}

		sess, _ := mgr.Store.Get(r, mgr.Config.CookieName)
		if sess == nil {
			return nil
		} // Safety Check

		// FIX PANIC: Cek Options sebelum akses
		if sess.Options == nil {
			sess.Options = &sessions.Options{}
		}

		sess.Options.MaxAge = -1
		sess.Options.Path = mgr.Config.Path
		sess.Options.Domain = mgr.Config.Domain

		return sess.Save(r, w)
	}

	// B. Stateless With Redis
	if mgr.Config.Driver == "stateless_with_redis" {
		tokenString := extractToken(r)
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(mgr.Config.SecretKey), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			sid := claims["sid"].(string)
			return utils.RedisDel(context.Background(), "sess:"+sid)
		}
	}

	return nil
}

// ==========================================
// INTERNAL HELPER
// ==========================================

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && strings.ToUpper(authHeader[0:6]) == "BEARER" {
		return authHeader[7:]
	}
	return ""
}

package login

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	resp "github.com/master-abror/zaframework/core/http"
	"github.com/master-abror/zaframework/core/session"
	"github.com/master-abror/zaframework/core/utils"
)

// Logout — POST /api/auth/logout
// Menghapus session dari Redis, clear cookie, invalidate token.
func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil token dari cookie
	tokenEncrypted := session.Get(r, "token")

	if tokenEncrypted != nil && tokenEncrypted != "" {
		tokenStr, ok := tokenEncrypted.(string)
		if ok && tokenStr != "" {
			// 2. Decrypt untuk dapat Redis key (UUID)
			decryptedKey, err := utils.Decrypt(tokenStr)
			if err == nil {
				// 3. Hapus JWT dari Redis
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				err = utils.RedisDel(ctx, string(decryptedKey))
				if err != nil {
					log.Printf("[LOGOUT] Redis delete error: %v", err)
				} else {
					log.Printf("[LOGOUT] Session deleted from Redis: %s", string(decryptedKey)[:8]+"...")
				}
			}
		}
	}

	// 4. Juga hapus session:<userID> jika ada
	if tokenEncrypted != nil && tokenEncrypted != "" {
		tokenStr, ok := tokenEncrypted.(string)
		if ok {
			decryptedKey, err := utils.Decrypt(tokenStr)
			if err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				// Get JWT to extract user_id
				jwtRaw, err := utils.RedisGet(ctx, string(decryptedKey))
				if err == nil && jwtRaw != "" {
					claims, err := utils.ValidateJWT(jwtRaw)
					if err == nil {
						if userID, ok := claims.User["id"].(string); ok {
							sessionKey := "session:" + userID
							_ = utils.RedisDel(ctx, sessionKey)
							log.Printf("[LOGOUT] Session key deleted: %s", sessionKey)
						}
					}
				}
			}
		}
	}

	// 5. Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "_authz",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   utils.GetEnv("SESSION_NAME", "za_session"),
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	fmt.Println("[LOGOUT] User logged out successfully")

	resp.ApiJSON(w, r, http.StatusOK, true, "Logout berhasil", nil)
}

package telegramlink

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Controller menangani HTTP request untuk link/unlink Telegram Chat ID.
type Controller struct {
	repo *Repository
}

// NewController membuat instance Controller baru.
func NewController() *Controller {
	return &Controller{repo: NewRepository()}
}

// LinkChatID menghubungkan Telegram Chat ID ke akun user.
//
// POST /api/telegram/link
// Header: X-User-ID: <user_id>
// Body: {"chat_id": "123456789"}
//
// Cara user mendapatkan Chat ID mereka:
// 1. User start/kirim pesan ke bot Telegram yang sudah dikonfigurasi
// 2. Bot membalas dengan chat_id mereka, atau user cek via @userinfobot
// 3. User input chat_id tersebut di halaman profil
func (ctrl *Controller) LinkChatID(c *gin.Context) {
	userID := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if userID == "" {
		userID = strings.TrimSpace(c.Query("user_id"))
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id wajib diisi (header X-User-ID)"})
		return
	}

	var req struct {
		ChatID string `json:"chat_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id wajib diisi"})
		return
	}

	chatID := strings.TrimSpace(req.ChatID)
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id tidak boleh kosong"})
		return
	}

	if err := ctrl.repo.SaveChatID(userID, chatID); err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Telegram berhasil dihubungkan. Kamu akan menerima notifikasi via Telegram."})
}

// UnlinkChatID melepas hubungan Telegram dari akun user.
//
// DELETE /api/telegram/unlink
// Header: X-User-ID: <user_id>
func (ctrl *Controller) UnlinkChatID(c *gin.Context) {
	userID := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if userID == "" {
		userID = strings.TrimSpace(c.Query("user_id"))
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id wajib diisi (header X-User-ID)"})
		return
	}

	if err := ctrl.repo.ClearChatID(userID); err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Telegram berhasil dilepas dari akun."})
}

// Status mengecek apakah user sudah menghubungkan Telegram.
//
// GET /api/telegram/status
// Header: X-User-ID: <user_id>
func (ctrl *Controller) Status(c *gin.Context) {
	userID := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if userID == "" {
		userID = strings.TrimSpace(c.Query("user_id"))
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id wajib diisi"})
		return
	}

	chatID, err := ctrl.repo.GetChatID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	linked := chatID != ""
	c.JSON(http.StatusOK, gin.H{
		"linked":  linked,
		"chat_id": chatID,
	})
}

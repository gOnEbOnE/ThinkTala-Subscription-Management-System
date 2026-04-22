package notification

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Controller menangani HTTP request untuk broadcast notifications.
// Controller hanya bicara dengan Service, tidak boleh langsung ke Repository/DB.
type Controller struct {
	svc *Service
}

// NewController membuat instance Controller baru.
func NewController() *Controller {
	return &Controller{svc: NewService()}
}

// List mengembalikan semua notification (admin/ops).
func (ctrl *Controller) List(c *gin.Context) {
	typeFilter := c.Query("type")
	statusFilter := c.Query("status")

	list, err := ctrl.svc.List(typeFilter, statusFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ListPublic mengembalikan notification aktif untuk role tertentu (dipakai frontend client).
func (ctrl *Controller) ListPublic(c *gin.Context) {
	role := strings.TrimSpace(c.GetHeader("X-User-Role"))
	if role == "" {
		role = strings.TrimSpace(c.Query("role"))
	}
	userID := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if userID == "" {
		userID = strings.TrimSpace(c.Query("user_id"))
	}

	list, err := ctrl.svc.ListPublic(role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Get mengembalikan satu notification berdasarkan ID.
func (ctrl *Controller) Get(c *gin.Context) {
	n, err := ctrl.svc.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification tidak ditemukan."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": n})
}

// Create membuat notification baru.
func (ctrl *Controller) Create(c *gin.Context) {
	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload tidak valid"})
		return
	}
	id := uuid.New().String()
	if err := ctrl.svc.Create(req, id); err != nil {
		if isValidationErr(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}, "message": "Notifikasi berhasil dibuat"})
}

// Update memperbarui notification berdasarkan ID.
func (ctrl *Controller) Update(c *gin.Context) {
	var req UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload tidak valid"})
		return
	}
	if err := ctrl.svc.Update(c.Param("id"), req); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification tidak ditemukan"})
			return
		}
		if isValidationErr(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil diperbarui"})
}

// Delete menghapus notification berdasarkan ID.
func (ctrl *Controller) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(c.Param("id")); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil dihapus"})
}

func isValidationErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "wajib") ||
		strings.Contains(msg, "tidak valid") ||
		strings.Contains(msg, "format expiry_date")
}

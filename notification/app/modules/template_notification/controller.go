package template

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller menangani HTTP request untuk notification templates.
// Controller hanya bicara dengan Service, tidak boleh langsung ke Repository/DB.
type Controller struct {
	svc *Service
}

// NewController membuat instance Controller baru dengan service yang diberikan.
func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// List mengembalikan semua template, bisa difilter via query (?channel=email&event_type=otp_verification).
func (ctrl *Controller) List(c *gin.Context) {
	list, err := ctrl.svc.List(c.Query("channel"), c.Query("event_type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Get mengembalikan satu template berdasarkan ID.
func (ctrl *Controller) Get(c *gin.Context) {
	t, err := ctrl.svc.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template tidak ditemukan."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

// Create membuat template baru.
func (ctrl *Controller) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	if err := ctrl.svc.Create(req, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}, "message": "Template berhasil dibuat."})
}

// Update memperbarui template berdasarkan ID.
func (ctrl *Controller) Update(c *gin.Context) {
	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.svc.Update(c.Param("id"), req); err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template berhasil diperbarui."})
}

// Delete menghapus template berdasarkan ID.
func (ctrl *Controller) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Template berhasil dihapus."})
}

// Send menerima event_type + channel + to + vars, cari template, render, dan kirim.
func (ctrl *Controller) Send(c *gin.Context) {
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.svc.Send(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi berhasil dikirim."})
}

// EventTypes mengembalikan daftar event_type yang sudah pernah di-dispatch.
func (ctrl *Controller) EventTypes(c *gin.Context) {
	list, err := ctrl.svc.ListEventTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// Logs mengembalikan log pengiriman notifikasi untuk monitoring.
// Query params: ?status=sent|failed|pending&channel=email|whatsapp|telegram&start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&search=...&limit=50&offset=0
func (ctrl *Controller) Logs(c *gin.Context) {
	status := strings.ToLower(strings.TrimSpace(c.Query("status")))
	if status == "all" {
		status = ""
	}
	channel := strings.ToLower(strings.TrimSpace(c.Query("channel")))
	if channel == "all" {
		channel = ""
	}
	search := strings.TrimSpace(c.Query("search"))
	if search != "" {
		search = "%" + search + "%"
	}

	var startDate *time.Time
	var endDate *time.Time
	if raw := strings.TrimSpace(c.Query("start_date")); raw != "" {
		parsed, err := time.Parse("2006-01-02", raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date tidak valid (YYYY-MM-DD)"})
			return
		}
		startDate = &parsed
	}
	if raw := strings.TrimSpace(c.Query("end_date")); raw != "" {
		parsed, err := time.Parse("2006-01-02", raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_date tidak valid (YYYY-MM-DD)"})
			return
		}
		parsed = parsed.AddDate(0, 0, 1)
		endDate = &parsed
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	list, total, err := ctrl.svc.GetLogs(status, channel, search, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   list,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// LogDetail mengembalikan detail log pengiriman notifikasi.
func (ctrl *Controller) LogDetail(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id wajib diisi"})
		return
	}
	logItem, err := ctrl.svc.GetLogDetail(id)
	if err != nil {
		log.Printf("[NOTIF LOG DETAIL] Error fetching log %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Log tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logItem})
}

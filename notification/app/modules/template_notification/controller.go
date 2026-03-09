package template

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller menangani HTTP request untuk notification templates.
// Controller hanya bicara dengan Service, tidak boleh langsung ke Repository/DB.
type Controller struct {
	svc *Service
}

// NewController membuat instance Controller baru.
func NewController() *Controller {
	return &Controller{svc: NewService()}
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

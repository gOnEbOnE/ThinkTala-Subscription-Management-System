package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	list, err := ctrl.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ListPublic mengembalikan notification aktif untuk role tertentu (dipakai frontend client).
func (ctrl *Controller) ListPublic(c *gin.Context) {
	list, err := ctrl.svc.ListPublic(c.Query("role"))
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	if err := ctrl.svc.Create(req, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": id}, "message": "Notification berhasil dibuat."})
}

// Update memperbarui notification berdasarkan ID.
func (ctrl *Controller) Update(c *gin.Context) {
	var req UpdateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.svc.Update(c.Param("id"), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification berhasil diperbarui."})
}

// Delete menghapus notification berdasarkan ID.
func (ctrl *Controller) Delete(c *gin.Context) {
	if err := ctrl.svc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification berhasil dihapus."})
}

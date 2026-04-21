package dashboard

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDashboardCustomers(c *gin.Context) {
	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	page, limit := normalizePageLimit(
		parseIntDefault(c.Query("page"), 1),
		parseIntDefault(c.Query("limit"), 10),
	)

	payload, err := h.service.GetDashboardCustomers(c.Request.Context(), dashboardQuery{
		PeriodType: periodType,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       page,
		Limit:      limit,
		Search:     strings.TrimSpace(c.Query("search")),
	})
	if err != nil {
		handleServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard customer berhasil dimuat",
		"data":    payload,
	})
}

func (h *Handler) GetDashboardCustomerDetail(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		setAuditError(c, "customer_id wajib diisi")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "customer_id wajib diisi"})
		return
	}

	payload, ok, err := h.service.GetCustomerDetail(c.Request.Context(), id)
	if err != nil {
		handleServerError(c, err)
		return
	}
	if !ok {
		setAuditError(c, "detail customer tidak tersedia")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Detail customer tidak tersedia."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func (h *Handler) GetDashboardPackages(c *gin.Context) {
	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	page, limit := normalizePageLimit(
		parseIntDefault(c.Query("page"), 1),
		parseIntDefault(c.Query("limit"), 10),
	)

	payload, err := h.service.GetDashboardPackages(c.Request.Context(), dashboardQuery{
		PeriodType: periodType,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       page,
		Limit:      limit,
		Search:     strings.TrimSpace(c.Query("search")),
	})
	if err != nil {
		handleServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard package berhasil dimuat",
		"data":    payload,
	})
}

func (h *Handler) GetDashboardPackageDetail(c *gin.Context) {
	packageID := strings.TrimSpace(c.Param("id"))
	if packageID == "" {
		setAuditError(c, "package_id wajib diisi")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "package_id wajib diisi"})
		return
	}

	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	payload, ok, err := h.service.GetPackageDetail(c.Request.Context(), packageID, dashboardQuery{
		PeriodType: periodType,
		StartDate:  startDate,
		EndDate:    endDate,
	})
	if err != nil {
		handleServerError(c, err)
		return
	}
	if !ok {
		setAuditError(c, "detail package tidak tersedia")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Detail package tidak tersedia."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func handleServerError(c *gin.Context, err error) {
	setAuditError(c, err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Terjadi kesalahan internal server",
	})
}

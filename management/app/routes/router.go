package routes

import (
	"net/http"

	"management/app/modules/dashboard"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, dashboardHandler *dashboard.Handler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "management"})
	})

	dashboardGroup := r.Group("/api/dashboard")
	dashboardGroup.Use(dashboard.RequireManagementRole())
	{
		dashboardGroup.GET("/customers", dashboardHandler.GetDashboardCustomers)
		dashboardGroup.GET("/customer/:id", dashboardHandler.GetDashboardCustomerDetail)
		dashboardGroup.GET("/packages", dashboard.RequirePackageDashboardRole(), dashboardHandler.GetDashboardPackages)
		dashboardGroup.GET("/package/:id", dashboard.RequirePackageDashboardRole(), dashboardHandler.GetDashboardPackageDetail)
	}
}

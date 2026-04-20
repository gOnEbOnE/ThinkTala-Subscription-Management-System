package dashboard

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireManagementRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := strings.ToUpper(strings.TrimSpace(c.GetHeader("X-User-Role")))
		if role == "MANAGEMENT" || role == "SUPERADMIN" || role == "ADMIN" {
			c.Next()
			return
		}

		setAuditError(c, "forbidden role")
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		c.Abort()
	}
}

func RequirePackageDashboardRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := strings.ToUpper(strings.TrimSpace(c.GetHeader("X-User-Role")))
		if role == "MANAGEMENT" || role == "ADMIN" || role == "SUPERADMIN" {
			c.Next()
			return
		}

		setAuditError(c, "forbidden role for package dashboard")
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
		})
		c.Abort()
	}
}

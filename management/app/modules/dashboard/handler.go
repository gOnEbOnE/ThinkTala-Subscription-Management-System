package dashboard

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

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

func (h *Handler) ExportDashboardPackages(c *gin.Context) {
	exportType := strings.TrimSpace(c.Query("type"))
	if exportType == "" {
		exportType = "summary"
	}

	if exportType != "summary" && exportType != "raw" {
		setAuditError(c, "invalid export type")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid export type. Use 'summary' or 'raw'"})
		return
	}

	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	payload, err := h.service.GetDashboardPackages(c.Request.Context(), dashboardQuery{
		PeriodType: periodType,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       1,
		Limit:      10000, // Get all records for export
		Search:     "",
	})
	if err != nil {
		handleServerError(c, err)
		return
	}

	var csvContent string
	if exportType == "summary" {
		csvContent = buildPackagesSummaryCSV(payload, startDate, endDate)
	} else {
		csvContent = buildPackagesRawCSV(payload)
	}

	filename := fmt.Sprintf("packages-export-%s.csv", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "text/csv;charset=utf-8")
	c.String(http.StatusOK, "\uFEFF"+csvContent)
}

func (h *Handler) ExportDashboardCustomers(c *gin.Context) {
	exportType := strings.TrimSpace(c.Query("type"))
	if exportType == "" {
		exportType = "summary"
	}

	if exportType != "summary" && exportType != "raw" {
		setAuditError(c, "invalid export type")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid export type. Use 'summary' or 'raw'"})
		return
	}

	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	payload, err := h.service.GetDashboardCustomers(c.Request.Context(), dashboardQuery{
		PeriodType: periodType,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       1,
		Limit:      10000, // Get all records for export
		Search:     "",
	})
	if err != nil {
		handleServerError(c, err)
		return
	}

	var csvContent string
	if exportType == "summary" {
		csvContent = buildCustomersSummaryCSV(payload)
	} else {
		csvContent = buildCustomersRawCSV(payload)
	}

	filename := fmt.Sprintf("customers-export-%s.csv", time.Now().Format("20060102-150405"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "text/csv;charset=utf-8")
	c.String(http.StatusOK, "\uFEFF"+csvContent)
}

func buildPackagesSummaryCSV(payload dashboardPackagePayload, startDate, endDate time.Time) string {
	var buf bytes.Buffer

	buf.WriteString("PACKAGE SALES DASHBOARD SUMMARY\r\n")
	buf.WriteString("Period," + startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02") + "\r\n")
	buf.WriteString("\r\n")

	buf.WriteString("SUMMARY METRICS\r\n")
	buf.WriteString("Metric,Value\r\n")
	buf.WriteString("Most Sold Package," + csvEscapeValue(payload.Summary.MostSoldPackage.PackageName) + "\r\n")
	buf.WriteString("Most Sold Transactions," + fmt.Sprintf("%d", payload.Summary.MostSoldPackage.TotalSales) + "\r\n")
	buf.WriteString("Highest Revenue Package," + csvEscapeValue(payload.Summary.HighestRevenuePackage.PackageName) + "\r\n")
	buf.WriteString("Highest Revenue," + fmt.Sprintf("%.2f", payload.Summary.HighestRevenuePackage.TotalRevenue) + "\r\n")
	buf.WriteString("Fastest Growth Package," + csvEscapeValue(payload.Summary.FastestGrowthPackage.PackageName) + "\r\n")
	buf.WriteString("Fastest Growth (%)," + fmt.Sprintf("%.2f", payload.Summary.FastestGrowthPackage.GrowthRate) + "\r\n")
	buf.WriteString("Total Revenue All Package," + fmt.Sprintf("%.2f", payload.Summary.TotalRevenueAllPackage) + "\r\n")
	buf.WriteString("\r\n")

	buf.WriteString("PACKAGE PERFORMANCE DETAILS\r\n")
	buf.WriteString("Package Name,Transactions,Revenue,Contribution (%),Growth (%)\r\n")
	if len(payload.Packages.Items) == 0 {
		buf.WriteString("Data tidak tersedia pada periode ini,,,\r\n")
	} else {
		for _, item := range payload.Packages.Items {
			buf.WriteString(csvEscapeValue(item.PackageName) + ",")
			buf.WriteString(fmt.Sprintf("%d", item.TotalTransactions) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.TotalRevenue) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.PercentageContribution) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.GrowthRate) + "\r\n")
		}
	}
	buf.WriteString("\r\n")

	buf.WriteString("REVENUE CONTRIBUTION\r\n")
	buf.WriteString("Package Name,Revenue,Contribution (%)\r\n")
	if len(payload.Charts.RevenueContribution) == 0 {
		buf.WriteString("Data tidak tersedia pada periode ini,,\r\n")
	} else {
		for _, item := range payload.Charts.RevenueContribution {
			buf.WriteString(csvEscapeValue(item.PackageName) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.Revenue) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.PercentageContribution) + "\r\n")
		}
	}

	return buf.String()
}

func buildPackagesRawCSV(payload dashboardPackagePayload) string {
	var buf bytes.Buffer

	buf.WriteString("PACKAGE SALES - RAW DATA\r\n")
	buf.WriteString("Exported at," + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
	buf.WriteString("\r\n")

	buf.WriteString("Package Name,Transactions,Revenue,Contribution (%),Growth (%)\r\n")
	if len(payload.Packages.Items) == 0 {
		buf.WriteString("No transactions found on this period\r\n")
	} else {
		for _, item := range payload.Packages.Items {
			buf.WriteString(csvEscapeValue(item.PackageName) + ",")
			buf.WriteString(fmt.Sprintf("%d", item.TotalTransactions) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.TotalRevenue) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.PercentageContribution) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.GrowthRate) + "\r\n")
		}
	}

	return buf.String()
}

func buildCustomersSummaryCSV(payload dashboardCustomerPayload) string {
	var buf bytes.Buffer

	buf.WriteString("CUSTOMER RETENTION & CHURN - SUMMARY\r\n")
	buf.WriteString("Period," + payload.Period.StartDate + " to " + payload.Period.EndDate + "\r\n")
	buf.WriteString("\r\n")

	buf.WriteString("SUMMARY METRICS\r\n")
	buf.WriteString("Metric,Value\r\n")
	buf.WriteString("Total Active Customers," + fmt.Sprintf("%d", payload.Summary.TotalActiveCustomers) + "\r\n")
	buf.WriteString("Total Loyal Customers," + fmt.Sprintf("%d", payload.Summary.TotalLoyalCustomers) + "\r\n")
	buf.WriteString("Total Churned Customers," + fmt.Sprintf("%d", payload.Summary.TotalChurnedCustomers) + "\r\n")
	buf.WriteString("Churn Rate (%)," + fmt.Sprintf("%.2f", payload.Summary.ChurnRate) + "\r\n")
	buf.WriteString("\r\n")

	buf.WriteString("MONTHLY NEW CUSTOMERS\r\n")
	buf.WriteString("Month,New Customers\r\n")
	for _, point := range payload.Charts.MonthlyNewCustomers {
		buf.WriteString(csvEscapeValue(point.Month) + ",")
		buf.WriteString(fmt.Sprintf("%d", point.Total) + "\r\n")
	}
	buf.WriteString("\r\n")

	buf.WriteString("MONTHLY CHURNED CUSTOMERS\r\n")
	buf.WriteString("Month,Churned Customers\r\n")
	for _, point := range payload.Charts.MonthlyChurnCustomers {
		buf.WriteString(csvEscapeValue(point.Month) + ",")
		buf.WriteString(fmt.Sprintf("%d", point.Total) + "\r\n")
	}
	buf.WriteString("\r\n")

	buf.WriteString("CHURN RATE TREND\r\n")
	buf.WriteString("Month,Churn Rate (%)\r\n")
	for _, point := range payload.Charts.ChurnRateTrend {
		buf.WriteString(csvEscapeValue(point.Month) + ",")
		buf.WriteString(fmt.Sprintf("%.2f", point.Rate) + "\r\n")
	}

	return buf.String()
}

func buildCustomersRawCSV(payload dashboardCustomerPayload) string {
	var buf bytes.Buffer

	buf.WriteString("TOP LOYAL CUSTOMERS - RAW DATA\r\n")
	buf.WriteString("Exported at," + time.Now().Format("2006-01-02 15:04:05") + "\r\n")
	buf.WriteString("\r\n")

	buf.WriteString("Customer Name,Email,Duration (months),Transactions,Total Spent,Last Active,Status\r\n")
	if len(payload.TopLoyalCustomers.Items) == 0 {
		buf.WriteString("No data available\r\n")
	} else {
		for _, item := range payload.TopLoyalCustomers.Items {
			buf.WriteString(csvEscapeValue(item.CustomerName) + ",")
			buf.WriteString(csvEscapeValue(item.Email) + ",")
			buf.WriteString(fmt.Sprintf("%d", item.Duration) + ",")
			buf.WriteString(fmt.Sprintf("%d", item.Transactions) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.TotalSpent) + ",")
			buf.WriteString(csvEscapeValue(item.LastActive) + ",")
			buf.WriteString("Loyal\r\n")
		}
	}
	buf.WriteString("\r\n")

	buf.WriteString("RECENTLY CHURNED CUSTOMERS - RAW DATA\r\n")
	buf.WriteString("Customer Name,Email,Last Subscription,Churn Date,Lifetime Value,Status\r\n")
	if len(payload.RecentlyChurnedCustomers) == 0 {
		buf.WriteString("No data available\r\n")
	} else {
		for _, item := range payload.RecentlyChurnedCustomers {
			buf.WriteString(csvEscapeValue(item.CustomerName) + ",")
			buf.WriteString(csvEscapeValue(item.Email) + ",")
			buf.WriteString(csvEscapeValue(item.LastSubscription) + ",")
			buf.WriteString(csvEscapeValue(item.ChurnDate) + ",")
			buf.WriteString(fmt.Sprintf("%.2f", item.LifetimeValue) + ",")
			buf.WriteString("Churned\r\n")
		}
	}

	return buf.String()
}

func csvEscapeValue(value string) string {
	if strings.ContainsAny(value, "\",\n\r") {
		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return value
}

func handleServerError(c *gin.Context, err error) {
	setAuditError(c, err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Terjadi kesalahan internal server",
	})
}

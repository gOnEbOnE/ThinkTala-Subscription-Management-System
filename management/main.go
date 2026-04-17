package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

type monthlyPoint struct {
	Month string `json:"month"`
	Total int    `json:"total"`
}

type churnRatePoint struct {
	Month string  `json:"month"`
	Rate  float64 `json:"rate"`
}

type loyalCustomerItem struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Email        string  `json:"email"`
	Duration     int     `json:"duration"`
	Transactions int     `json:"transactions"`
	TotalSpent   float64 `json:"total_spent"`
	LastActive   string  `json:"last_active"`
}

type churnedCustomerItem struct {
	CustomerID       string  `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	Email            string  `json:"email"`
	LastSubscription string  `json:"last_subscription"`
	ChurnDate        string  `json:"churn_date"`
	LifetimeValue    float64 `json:"lifetime_value"`
}

type dashboardCustomerPayload struct {
	Period struct {
		Type      string `json:"type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"period"`
	Summary struct {
		TotalActiveCustomers      int     `json:"total_active_customers"`
		TotalLoyalCustomers       int     `json:"total_loyal_customers"`
		TotalChurnedCustomers     int     `json:"total_churned_customers"`
		ChurnRate                 float64 `json:"churn_rate"`
		ActiveCustomersChangePct  float64 `json:"active_customers_change_pct"`
		ChurnRateChangePct        float64 `json:"churn_rate_change_pct"`
		ActiveCustomersLastPeriod int     `json:"active_customers_last_period"`
		ChurnRateLastPeriod       float64 `json:"churn_rate_last_period"`
	} `json:"summary"`
	Charts struct {
		MonthlyNewCustomers   []monthlyPoint   `json:"monthly_new_customers"`
		MonthlyChurnCustomers []monthlyPoint   `json:"monthly_churn_customers"`
		ChurnRateTrend        []churnRatePoint `json:"churn_rate_trend"`
	} `json:"charts"`
	TopLoyalCustomers struct {
		Items      []loyalCustomerItem `json:"items"`
		Pagination struct {
			Page       int `json:"page"`
			Limit      int `json:"limit"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"pagination"`
	} `json:"top_loyal_customers"`
	RecentlyChurnedCustomers []churnedCustomerItem `json:"recently_churned_customers"`
}

type customerDetailPayload struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	Email        string  `json:"email"`
	Duration     int     `json:"duration"`
	Transactions int     `json:"transactions"`
	TotalSpent   float64 `json:"total_spent"`
	LastActive   string  `json:"last_active"`
}

type requestAudit struct {
	Endpoint string
	Method   string
	Role     string
	Status   int
	Message  string
	Payload  string
}

type dummyCustomer struct {
	CustomerID       string
	CustomerName     string
	Email            string
	JoinDate         time.Time
	ActiveUntil      time.Time
	LastActive       time.Time
	Transactions     int
	TotalSpent       float64
	LastSubscription string
	ChurnDate        *time.Time
}

type packageSalesDistributionPoint struct {
	PackageID   string `json:"package_id"`
	PackageName string `json:"package_name"`
	TotalSales  int    `json:"total_sales"`
}

type packageRevenueContributionPoint struct {
	PackageID              string  `json:"package_id"`
	PackageName            string  `json:"package_name"`
	Revenue                float64 `json:"revenue"`
	PercentageContribution float64 `json:"percentage_contribution"`
}

type packageSalesTrendPoint struct {
	Period     string `json:"period"`
	Premium    int    `json:"premium"`
	Enterprise int    `json:"enterprise"`
	Starter    int    `json:"starter"`
}

type packagePerformanceItem struct {
	PackageID              string  `json:"package_id"`
	PackageName            string  `json:"package_name"`
	TotalTransactions      int     `json:"total_transactions"`
	TotalRevenue           float64 `json:"total_revenue"`
	PercentageContribution float64 `json:"percentage_contribution"`
	GrowthRate             float64 `json:"growth_rate"`
}

type dashboardPackagePayload struct {
	Period struct {
		Type      string `json:"type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	} `json:"period"`
	Summary struct {
		MostSoldPackage struct {
			PackageID   string `json:"package_id"`
			PackageName string `json:"package_name"`
			TotalSales  int    `json:"total_sales"`
		} `json:"most_sold_package"`
		HighestRevenuePackage struct {
			PackageID    string  `json:"package_id"`
			PackageName  string  `json:"package_name"`
			TotalRevenue float64 `json:"total_revenue"`
		} `json:"highest_revenue_package"`
		FastestGrowthPackage struct {
			PackageID   string  `json:"package_id"`
			PackageName string  `json:"package_name"`
			GrowthRate  float64 `json:"growth_rate"`
		} `json:"fastest_growth_package"`
		TotalRevenueAllPackage float64 `json:"total_revenue_all_package"`
		PeriodLabel            string  `json:"period_label"`
	} `json:"summary"`
	Charts struct {
		PackageSalesDistribution []packageSalesDistributionPoint   `json:"package_sales_distribution"`
		RevenueContribution      []packageRevenueContributionPoint `json:"revenue_contribution"`
		SalesTrend               []packageSalesTrendPoint          `json:"sales_trend"`
	} `json:"charts"`
	Packages struct {
		Items      []packagePerformanceItem `json:"items"`
		Pagination struct {
			Page       int `json:"page"`
			Limit      int `json:"limit"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"pagination"`
	} `json:"packages"`
}

type packageDetailTrendPoint struct {
	Period string `json:"period"`
	Sales  int    `json:"sales"`
}

type packageDetailPayload struct {
	PackageID              string                    `json:"package_id"`
	PackageName            string                    `json:"package_name"`
	TotalTransactions      int                       `json:"total_transactions"`
	TotalRevenue           float64                   `json:"total_revenue"`
	PercentageContribution float64                   `json:"percentage_contribution"`
	GrowthRate             float64                   `json:"growth_rate"`
	Period                 map[string]string         `json:"period"`
	Trend                  []packageDetailTrendPoint `json:"trend"`
}

type packageCatalogItem struct {
	ID        string
	Name      string
	ShortName string
	Price     float64
}

type dummyPackageTransaction struct {
	PackageID   string
	PackageName string
	Price       float64
	Status      string
	CreatedAt   time.Time
}

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "thinknalyze"),
		getEnv("DB_SSLMODE", "disable"),
	)

	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("[MANAGEMENT] gagal koneksi database: %v", err)
	}
	defer db.Close()

	if err := ensureManagementSchema(context.Background()); err != nil {
		log.Fatalf("[MANAGEMENT] gagal menyiapkan schema: %v", err)
	}

	r := gin.Default()
	r.Use(auditMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "management"})
	})

	dashboard := r.Group("/api/dashboard")
	dashboard.Use(requireManagementRole())
	{
		dashboard.GET("/customers", getDashboardCustomers)
		dashboard.GET("/customer/:id", getDashboardCustomerDetail)
		dashboard.GET("/packages", requirePackageDashboardRole(), getDashboardPackages)
		dashboard.GET("/package/:id", requirePackageDashboardRole(), getDashboardPackageDetail)
	}

	port := getEnv("PORT", "5006")
	log.Printf("[MANAGEMENT] service running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[MANAGEMENT] gagal menjalankan service: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireManagementRole() gin.HandlerFunc {
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

func requirePackageDashboardRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := strings.ToUpper(strings.TrimSpace(c.GetHeader("X-User-Role")))
		if role == "MANAGEMENT" || role == "ADMIN" {
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

func getDashboardCustomers(c *gin.Context) {
	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	page := parseIntDefault(c.Query("page"), 1)
	if page < 1 {
		page = 1
	}
	limit := parseIntDefault(c.Query("limit"), 10)
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	search := strings.TrimSpace(c.Query("search"))

	payload := buildDummyDashboardPayload(periodType, startDate, endDate, page, limit, search)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard customer berhasil dimuat",
		"data":    payload,
	})
}

func getDashboardCustomerDetail(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "customer_id wajib diisi"})
		return
	}

	payload, ok := buildDummyCustomerDetail(id, time.Now())
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Detail customer tidak tersedia."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func getDashboardPackages(c *gin.Context) {
	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	page := parseIntDefault(c.Query("page"), 1)
	if page < 1 {
		page = 1
	}
	limit := parseIntDefault(c.Query("limit"), 10)
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	search := strings.TrimSpace(c.Query("search"))
	payload := buildDummyPackageDashboardPayload(periodType, startDate, endDate, page, limit, search)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard package berhasil dimuat",
		"data":    payload,
	})
}

func getDashboardPackageDetail(c *gin.Context) {
	packageID := strings.TrimSpace(c.Param("id"))
	if packageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "package_id wajib diisi"})
		return
	}

	periodType, startDate, endDate, err := parsePeriod(c)
	if err != nil {
		setAuditError(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	payload, ok := buildDummyPackageDetailPayload(packageID, periodType, startDate, endDate)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Detail package tidak tersedia."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func buildDummyPackageDashboardPayload(periodType string, startDate, endDate time.Time, page, limit int, search string) dashboardPackagePayload {
	transactions := getDummyPackageTransactions(startDate.Location())
	catalog := getPackageCatalog()
	catalogMap := make(map[string]packageCatalogItem, len(catalog))
	for _, item := range catalog {
		catalogMap[item.ID] = item
	}

	currentMetrics := aggregatePackageMetrics(transactions, startDate, endDate, catalogMap)
	prevStart, prevEnd := previousPeriod(startDate, endDate)
	prevMetrics := aggregatePackageMetrics(transactions, prevStart, prevEnd, catalogMap)

	payload := dashboardPackagePayload{}
	payload.Period.Type = periodType
	payload.Period.StartDate = startDate.Format("2006-01-02")
	payload.Period.EndDate = endDate.Format("2006-01-02")
	payload.Summary.PeriodLabel = packagePeriodLabel(periodType)

	totalRevenueAll := 0.0
	totalTransactionsAll := 0
	for _, metric := range currentMetrics {
		totalRevenueAll += metric.TotalRevenue
		totalTransactionsAll += metric.TotalTransactions
	}

	if totalTransactionsAll == 0 {
		payload.Summary.MostSoldPackage.PackageName = "-"
		payload.Summary.HighestRevenuePackage.PackageName = "-"
		payload.Summary.FastestGrowthPackage.PackageName = "-"
		payload.Packages.Items = []packagePerformanceItem{}
		payload.Charts.PackageSalesDistribution = []packageSalesDistributionPoint{}
		payload.Charts.RevenueContribution = []packageRevenueContributionPoint{}
		payload.Charts.SalesTrend = []packageSalesTrendPoint{}
		payload.Packages.Pagination.Page = page
		payload.Packages.Pagination.Limit = limit
		payload.Packages.Pagination.Total = 0
		payload.Packages.Pagination.TotalPages = 0
		return payload
	}

	metricsList := make([]packagePerformanceItem, 0, len(currentMetrics))
	mostSold := packagePerformanceItem{PackageName: "-"}
	highestRevenue := packagePerformanceItem{PackageName: "-"}
	fastestGrowth := packagePerformanceItem{PackageName: "-", GrowthRate: -999999}

	for _, item := range catalog {
		metric, ok := currentMetrics[item.ID]
		if !ok || metric.TotalTransactions == 0 {
			continue
		}

		prevTx := 0
		if prev, hasPrev := prevMetrics[item.ID]; hasPrev {
			prevTx = prev.TotalTransactions
		}
		growth := compareRate(metric.TotalTransactions, prevTx)
		contribution := 0.0
		if totalRevenueAll > 0 {
			contribution = round2((metric.TotalRevenue / totalRevenueAll) * 100)
		}

		itemMetric := packagePerformanceItem{
			PackageID:              metric.PackageID,
			PackageName:            metric.PackageName,
			TotalTransactions:      metric.TotalTransactions,
			TotalRevenue:           metric.TotalRevenue,
			PercentageContribution: contribution,
			GrowthRate:             growth,
		}
		metricsList = append(metricsList, itemMetric)

		if itemMetric.TotalTransactions > mostSold.TotalTransactions {
			mostSold = itemMetric
		}
		if itemMetric.TotalRevenue > highestRevenue.TotalRevenue {
			highestRevenue = itemMetric
		}
		if itemMetric.GrowthRate > fastestGrowth.GrowthRate {
			fastestGrowth = itemMetric
		}
	}

	payload.Summary.MostSoldPackage.PackageID = mostSold.PackageID
	payload.Summary.MostSoldPackage.PackageName = mostSold.PackageName
	payload.Summary.MostSoldPackage.TotalSales = mostSold.TotalTransactions

	payload.Summary.HighestRevenuePackage.PackageID = highestRevenue.PackageID
	payload.Summary.HighestRevenuePackage.PackageName = highestRevenue.PackageName
	payload.Summary.HighestRevenuePackage.TotalRevenue = highestRevenue.TotalRevenue

	payload.Summary.FastestGrowthPackage.PackageID = fastestGrowth.PackageID
	payload.Summary.FastestGrowthPackage.PackageName = fastestGrowth.PackageName
	payload.Summary.FastestGrowthPackage.GrowthRate = fastestGrowth.GrowthRate
	payload.Summary.TotalRevenueAllPackage = totalRevenueAll

	searchLower := strings.ToLower(strings.TrimSpace(search))
	filteredItems := make([]packagePerformanceItem, 0, len(metricsList))
	for _, metric := range metricsList {
		if searchLower != "" && !strings.Contains(strings.ToLower(metric.PackageName), searchLower) {
			continue
		}
		filteredItems = append(filteredItems, metric)
	}

	sort.Slice(filteredItems, func(i, j int) bool {
		if filteredItems[i].TotalTransactions == filteredItems[j].TotalTransactions {
			return filteredItems[i].TotalRevenue > filteredItems[j].TotalRevenue
		}
		return filteredItems[i].TotalTransactions > filteredItems[j].TotalTransactions
	})

	totalItems := len(filteredItems)
	startIdx := (page - 1) * limit
	if startIdx > totalItems {
		startIdx = totalItems
	}
	endIdx := startIdx + limit
	if endIdx > totalItems {
		endIdx = totalItems
	}

	payload.Packages.Items = filteredItems[startIdx:endIdx]
	payload.Packages.Pagination.Page = page
	payload.Packages.Pagination.Limit = limit
	payload.Packages.Pagination.Total = totalItems
	payload.Packages.Pagination.TotalPages = calcTotalPages(totalItems, limit)

	distribution := make([]packageSalesDistributionPoint, 0, len(catalog))
	contribution := make([]packageRevenueContributionPoint, 0, len(catalog))
	for _, item := range catalog {
		metric := currentMetrics[item.ID]
		totalSales := 0
		revenue := 0.0
		percentage := 0.0
		if metric != nil {
			totalSales = metric.TotalTransactions
			revenue = metric.TotalRevenue
			if totalRevenueAll > 0 {
				percentage = round2((revenue / totalRevenueAll) * 100)
			}
		}

		distribution = append(distribution, packageSalesDistributionPoint{
			PackageID:   item.ID,
			PackageName: item.ShortName,
			TotalSales:  totalSales,
		})
		contribution = append(contribution, packageRevenueContributionPoint{
			PackageID:              item.ID,
			PackageName:            item.ShortName,
			Revenue:                revenue,
			PercentageContribution: percentage,
		})
	}
	payload.Charts.PackageSalesDistribution = distribution
	payload.Charts.RevenueContribution = contribution

	months := buildMonthBuckets(startDate, endDate)
	trend := make([]packageSalesTrendPoint, 0, len(months))
	for _, monthStart := range months {
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
		monthMetrics := aggregatePackageMetrics(transactions, monthStart, monthEnd, catalogMap)
		trend = append(trend, packageSalesTrendPoint{
			Period:     monthStart.Format("Jan 2006"),
			Premium:    getMetricTransactions(monthMetrics, "premium"),
			Enterprise: getMetricTransactions(monthMetrics, "enterprise"),
			Starter:    getMetricTransactions(monthMetrics, "starter"),
		})
	}
	payload.Charts.SalesTrend = trend

	return payload
}

func buildDummyPackageDetailPayload(packageID, periodType string, startDate, endDate time.Time) (packageDetailPayload, bool) {
	catalog := getPackageCatalog()
	catalogMap := make(map[string]packageCatalogItem, len(catalog))
	for _, item := range catalog {
		catalogMap[item.ID] = item
	}

	selected, ok := catalogMap[strings.ToLower(packageID)]
	if !ok {
		selected, ok = catalogMap[packageID]
	}
	if !ok {
		return packageDetailPayload{}, false
	}

	transactions := getDummyPackageTransactions(startDate.Location())
	currentMetrics := aggregatePackageMetrics(transactions, startDate, endDate, catalogMap)
	prevStart, prevEnd := previousPeriod(startDate, endDate)
	prevMetrics := aggregatePackageMetrics(transactions, prevStart, prevEnd, catalogMap)

	totalRevenueAll := 0.0
	for _, metric := range currentMetrics {
		totalRevenueAll += metric.TotalRevenue
	}

	current := currentMetrics[selected.ID]
	totalTransactions := 0
	totalRevenue := 0.0
	if current != nil {
		totalTransactions = current.TotalTransactions
		totalRevenue = current.TotalRevenue
	}

	prevTx := 0
	if prev := prevMetrics[selected.ID]; prev != nil {
		prevTx = prev.TotalTransactions
	}
	growth := compareRate(totalTransactions, prevTx)
	contribution := 0.0
	if totalRevenueAll > 0 {
		contribution = round2((totalRevenue / totalRevenueAll) * 100)
	}

	months := buildMonthBuckets(startDate, endDate)
	trend := make([]packageDetailTrendPoint, 0, len(months))
	for _, monthStart := range months {
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
		monthMetrics := aggregatePackageMetrics(transactions, monthStart, monthEnd, catalogMap)
		trend = append(trend, packageDetailTrendPoint{
			Period: monthStart.Format("Jan 2006"),
			Sales:  getMetricTransactions(monthMetrics, selected.ID),
		})
	}

	payload := packageDetailPayload{
		PackageID:              selected.ID,
		PackageName:            selected.Name,
		TotalTransactions:      totalTransactions,
		TotalRevenue:           totalRevenue,
		PercentageContribution: contribution,
		GrowthRate:             growth,
		Period: map[string]string{
			"type":       periodType,
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		Trend: trend,
	}

	return payload, true
}

func getPackageCatalog() []packageCatalogItem {
	return []packageCatalogItem{
		{ID: "starter", Name: "Starter Plan", ShortName: "Starter", Price: 120000},
		{ID: "premium", Name: "Premium Plan", ShortName: "Premium", Price: 249000},
		{ID: "enterprise", Name: "Enterprise Plan", ShortName: "Enterprise", Price: 399000},
		{ID: "custom", Name: "Custom Plan", ShortName: "Custom", Price: 589000},
	}
}

func getDummyPackageTransactions(loc *time.Location) []dummyPackageTransaction {
	if loc == nil {
		loc = time.Local
	}

	catalog := getPackageCatalog()
	catalogMap := make(map[string]packageCatalogItem, len(catalog))
	for _, item := range catalog {
		catalogMap[item.ID] = item
	}

	seeds := []struct {
		Year       int
		Month      time.Month
		Starter    int
		Premium    int
		Enterprise int
		Custom     int
	}{
		{Year: 2025, Month: time.November, Starter: 48, Premium: 72, Enterprise: 28, Custom: 13},
		{Year: 2025, Month: time.December, Starter: 52, Premium: 84, Enterprise: 30, Custom: 16},
		{Year: 2026, Month: time.January, Starter: 60, Premium: 96, Enterprise: 40, Custom: 20},
		{Year: 2026, Month: time.February, Starter: 72, Premium: 118, Enterprise: 52, Custom: 24},
		{Year: 2026, Month: time.March, Starter: 84, Premium: 136, Enterprise: 60, Custom: 30},
		{Year: 2026, Month: time.April, Starter: 110, Premium: 167, Enterprise: 74, Custom: 33},
		{Year: 2026, Month: time.May, Starter: 124, Premium: 188, Enterprise: 86, Custom: 39},
		{Year: 2026, Month: time.June, Starter: 138, Premium: 210, Enterprise: 95, Custom: 45},
		{Year: 2026, Month: time.July, Starter: 152, Premium: 232, Enterprise: 104, Custom: 52},
		{Year: 2026, Month: time.August, Starter: 168, Premium: 260, Enterprise: 116, Custom: 58},
	}

	all := make([]dummyPackageTransaction, 0)
	appendTx := func(seedYear int, seedMonth time.Month, packageID string, count int) {
		pkg := catalogMap[packageID]
		for i := 0; i < count; i++ {
			status := "SUCCESS"
			if i%4 == 0 {
				status = "PAID"
			}
			all = append(all, dummyPackageTransaction{
				PackageID:   pkg.ID,
				PackageName: pkg.Name,
				Price:       pkg.Price,
				Status:      status,
				CreatedAt:   time.Date(seedYear, seedMonth, (i%27)+1, 10+(i%8), 0, 0, 0, loc),
			})
		}

		failedCount := count / 20
		if failedCount < 1 {
			failedCount = 1
		}
		for i := 0; i < failedCount; i++ {
			all = append(all, dummyPackageTransaction{
				PackageID:   pkg.ID,
				PackageName: pkg.Name,
				Price:       pkg.Price,
				Status:      "FAILED",
				CreatedAt:   time.Date(seedYear, seedMonth, (i%27)+1, 18, 0, 0, 0, loc),
			})
		}
	}

	for _, seed := range seeds {
		appendTx(seed.Year, seed.Month, "starter", seed.Starter)
		appendTx(seed.Year, seed.Month, "premium", seed.Premium)
		appendTx(seed.Year, seed.Month, "enterprise", seed.Enterprise)
		appendTx(seed.Year, seed.Month, "custom", seed.Custom)
	}

	return all
}

func aggregatePackageMetrics(transactions []dummyPackageTransaction, startDate, endDate time.Time, catalogMap map[string]packageCatalogItem) map[string]*packagePerformanceItem {
	metrics := make(map[string]*packagePerformanceItem)
	for _, tx := range transactions {
		status := strings.ToUpper(strings.TrimSpace(tx.Status))
		if status != "SUCCESS" && status != "PAID" {
			continue
		}
		if !isWithinRange(tx.CreatedAt, startDate, endDate) {
			continue
		}

		meta := catalogMap[tx.PackageID]
		if _, ok := metrics[tx.PackageID]; !ok {
			metrics[tx.PackageID] = &packagePerformanceItem{
				PackageID:   tx.PackageID,
				PackageName: meta.Name,
			}
		}
		metrics[tx.PackageID].TotalTransactions++
		metrics[tx.PackageID].TotalRevenue += tx.Price
	}
	return metrics
}

func getMetricTransactions(metrics map[string]*packagePerformanceItem, packageID string) int {
	if metric, ok := metrics[packageID]; ok {
		return metric.TotalTransactions
	}
	return 0
}

func packagePeriodLabel(periodType string) string {
	switch strings.ToLower(strings.TrimSpace(periodType)) {
	case "yearly":
		return "This Year"
	case "custom":
		return "Custom Period"
	default:
		return "This Month"
	}
}

func parsePeriod(c *gin.Context) (string, time.Time, time.Time, error) {
	now := time.Now()
	loc := now.Location()

	monthRaw := strings.TrimSpace(c.Query("month"))
	yearRaw := strings.TrimSpace(c.Query("year"))
	startRaw := strings.TrimSpace(c.Query("start_date"))
	endRaw := strings.TrimSpace(c.Query("end_date"))

	if monthRaw != "" && yearRaw == "" {
		return "", time.Time{}, time.Time{}, errors.New("month wajib dengan year")
	}

	if startRaw != "" || endRaw != "" {
		if startRaw == "" || endRaw == "" {
			return "", time.Time{}, time.Time{}, errors.New("start_date dan end_date wajib diisi bersamaan")
		}

		startDate, err := time.ParseInLocation("2006-01-02", startRaw, loc)
		if err != nil || startDate.Format("2006-01-02") != startRaw {
			return "", time.Time{}, time.Time{}, errors.New("format start_date tidak valid (YYYY-MM-DD)")
		}

		endDate, err := time.ParseInLocation("2006-01-02", endRaw, loc)
		if err != nil || endDate.Format("2006-01-02") != endRaw {
			return "", time.Time{}, time.Time{}, errors.New("format end_date tidak valid (YYYY-MM-DD)")
		}

		if endDate.Before(startDate) {
			return "", time.Time{}, time.Time{}, errors.New("Rentang tanggal tidak valid.")
		}

		return "custom", startDate, endDate.Add(24*time.Hour - time.Nanosecond), nil
	}

	if monthRaw != "" || yearRaw != "" {
		year, err := strconv.Atoi(yearRaw)
		if err != nil || year < 1900 || year > 3000 {
			return "", time.Time{}, time.Time{}, errors.New("year tidak valid")
		}

		if monthRaw == "" {
			start := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)
			end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
			return "yearly", start, end, nil
		}

		month, convErr := strconv.Atoi(monthRaw)
		if convErr != nil || month < 1 || month > 12 {
			return "", time.Time{}, time.Time{}, errors.New("month tidak valid")
		}

		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		return "monthly", start, end, nil
	}

	period := strings.ToLower(strings.TrimSpace(c.Query("period")))
	if period == "yearly" {
		start := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, loc)
		return "yearly", start, start.AddDate(1, 0, 0).Add(-time.Nanosecond), nil
	}

	base := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	return "monthly", base, base.AddDate(0, 1, 0).Add(-time.Nanosecond), nil
}

func previousPeriod(start, end time.Time) (time.Time, time.Time) {
	duration := end.Sub(start) + time.Nanosecond
	prevEnd := start.Add(-time.Nanosecond)
	prevStart := prevEnd.Add(-duration + time.Nanosecond)
	return prevStart, prevEnd
}

func countActiveCustomers(ctx context.Context, startDate, endDate time.Time) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => COALESCE(p.duration, 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		)
		SELECT COUNT(DISTINCT user_id)
		FROM paid_orders
		WHERE created_at <= $2
		  AND expiry_at >= $1
	`
	var total int
	err := db.QueryRow(ctx, query, startDate, endDate).Scan(&total)
	return total, err
}

func countActiveAtDate(ctx context.Context, at time.Time) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => COALESCE(p.duration, 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		)
		SELECT COUNT(DISTINCT user_id)
		FROM paid_orders
		WHERE created_at <= $1
		  AND expiry_at > $1
	`
	var total int
	err := db.QueryRow(ctx, query, at).Scan(&total)
	return total, err
}

func countChurnedCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => COALESCE(p.duration, 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		),
		churn_events AS (
			SELECT DISTINCT po.user_id
			FROM paid_orders po
			WHERE po.expiry_at >= $1
			  AND po.expiry_at <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = po.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > po.expiry_at
				  AND r.created_at <= (po.expiry_at + make_interval(days => $3))
			)
		)
		SELECT COUNT(*) FROM churn_events
	`
	var total int
	err := db.QueryRow(ctx, query, startDate, endDate, graceDays).Scan(&total)
	return total, err
}

func countLoyalCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				COALESCE(p.duration, 1) AS duration_months,
				(o.created_at + make_interval(months => COALESCE(p.duration, 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		),
		churn_events AS (
			SELECT DISTINCT po.user_id
			FROM paid_orders po
			WHERE po.expiry_at >= $1
			  AND po.expiry_at <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = po.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > po.expiry_at
				  AND r.created_at <= (po.expiry_at + make_interval(days => $3))
			)
		),
		aggregate_subs AS (
			SELECT user_id, SUM(duration_months)::int AS total_months
			FROM paid_orders
			WHERE created_at <= $2
			GROUP BY user_id
		)
		SELECT COUNT(*)
		FROM aggregate_subs a
		WHERE a.total_months >= 12
		  AND NOT EXISTS (SELECT 1 FROM churn_events c WHERE c.user_id = a.user_id)
	`
	var total int
	err := db.QueryRow(ctx, query, startDate, endDate, graceDays).Scan(&total)
	return total, err
}

func fetchMonthlyNewCustomers(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
	query := `
		WITH first_paid AS (
			SELECT user_id, MIN(created_at) AS first_paid_at
			FROM subscription.orders
			WHERE status = 'PAID'
			GROUP BY user_id
		)
		SELECT date_trunc('month', first_paid_at) AS month_start, COUNT(*)::int AS total
		FROM first_paid
		WHERE first_paid_at >= $1
		  AND first_paid_at <= $2
		GROUP BY 1
		ORDER BY 1
	`
	rows, err := db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var monthStart time.Time
		var total int
		if scanErr := rows.Scan(&monthStart, &total); scanErr != nil {
			return nil, scanErr
		}
		result[monthKey(monthStart)] = total
	}
	return result, rows.Err()
}

func fetchMonthlyChurnCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (map[string]int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => COALESCE(p.duration, 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		),
		churn_events AS (
			SELECT po.user_id, po.expiry_at
			FROM paid_orders po
			WHERE po.expiry_at >= $1
			  AND po.expiry_at <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = po.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > po.expiry_at
				  AND r.created_at <= (po.expiry_at + make_interval(days => $3))
			)
		)
		SELECT date_trunc('month', expiry_at) AS month_start, COUNT(DISTINCT user_id)::int AS total
		FROM churn_events
		GROUP BY 1
		ORDER BY 1
	`
	rows, err := db.Query(ctx, query, startDate, endDate, graceDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var monthStart time.Time
		var total int
		if scanErr := rows.Scan(&monthStart, &total); scanErr != nil {
			return nil, scanErr
		}
		result[monthKey(monthStart)] = total
	}
	return result, rows.Err()
}

func fetchTopLoyalCustomers(ctx context.Context, endDate time.Time, search string, page, limit int) ([]loyalCustomerItem, int, error) {
	offset := (page - 1) * limit
	search = strings.ToLower(strings.TrimSpace(search))

	countQuery := `
		WITH user_agg AS (
			SELECT
				o.user_id,
				SUM(COALESCE(p.duration, 1))::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(o.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(o.created_at) AS last_active
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
			  AND o.created_at <= $1
			GROUP BY o.user_id
		)
		SELECT COUNT(*)
		FROM user_agg ua
		JOIN users u ON u.id = ua.user_id
		WHERE ua.duration_months >= 12
		  AND (
			$2 = ''
			OR LOWER(COALESCE(u.name, '')) LIKE '%' || $2 || '%'
			OR LOWER(COALESCE(u.email, '')) LIKE '%' || $2 || '%'
		  )
	`
	var total int
	if err := db.QueryRow(ctx, countQuery, endDate, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		WITH user_agg AS (
			SELECT
				o.user_id,
				SUM(COALESCE(p.duration, 1))::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(o.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(o.created_at) AS last_active
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
			  AND o.created_at <= $1
			GROUP BY o.user_id
		)
		SELECT
			u.id::text,
			COALESCE(u.name, 'Unknown') AS customer_name,
			COALESCE(u.email, '') AS email,
			ua.duration_months,
			ua.transactions,
			ua.total_spent,
			ua.last_active
		FROM user_agg ua
		JOIN users u ON u.id = ua.user_id
		WHERE ua.duration_months >= 12
		  AND (
			$2 = ''
			OR LOWER(COALESCE(u.name, '')) LIKE '%' || $2 || '%'
			OR LOWER(COALESCE(u.email, '')) LIKE '%' || $2 || '%'
		  )
		ORDER BY ua.duration_months DESC, ua.total_spent DESC, ua.last_active DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := db.Query(ctx, query, endDate, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]loyalCustomerItem, 0)
	for rows.Next() {
		var item loyalCustomerItem
		var lastActive time.Time
		if scanErr := rows.Scan(
			&item.CustomerID,
			&item.CustomerName,
			&item.Email,
			&item.Duration,
			&item.Transactions,
			&item.TotalSpent,
			&lastActive,
		); scanErr != nil {
			return nil, 0, scanErr
		}
		item.LastActive = lastActive.Format(time.RFC3339)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if items == nil {
		items = []loyalCustomerItem{}
	}
	return items, total, nil
}

func fetchRecentlyChurnedCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) ([]churnedCustomerItem, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				o.total_price,
				COALESCE(p.name, 'Unknown Plan') AS package_name,
				(o.created_at + make_interval(months => COALESCE(p.duration, 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		),
		churn_events AS (
			SELECT
				po.user_id,
				po.package_name,
				po.expiry_at,
				po.total_price
			FROM paid_orders po
			WHERE po.expiry_at >= $1
			  AND po.expiry_at <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = po.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > po.expiry_at
				  AND r.created_at <= (po.expiry_at + make_interval(days => $3))
			)
		),
		ranked AS (
			SELECT
				ce.user_id,
				ce.package_name,
				ce.expiry_at,
				ROW_NUMBER() OVER (PARTITION BY ce.user_id ORDER BY ce.expiry_at DESC) AS rn
			FROM churn_events ce
		),
		ltv AS (
			SELECT user_id, COALESCE(SUM(total_price), 0)::numeric(18,2) AS lifetime_value
			FROM subscription.orders
			WHERE status = 'PAID'
			GROUP BY user_id
		)
		SELECT
			u.id::text,
			COALESCE(u.name, 'Unknown') AS customer_name,
			COALESCE(u.email, '') AS email,
			r.package_name,
			r.expiry_at,
			COALESCE(l.lifetime_value, 0)
		FROM ranked r
		JOIN users u ON u.id = r.user_id
		LEFT JOIN ltv l ON l.user_id = r.user_id
		WHERE r.rn = 1
		ORDER BY r.expiry_at DESC
		LIMIT 10
	`

	rows, err := db.Query(ctx, query, startDate, endDate, graceDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]churnedCustomerItem, 0)
	for rows.Next() {
		var item churnedCustomerItem
		var churnDate time.Time
		if scanErr := rows.Scan(
			&item.CustomerID,
			&item.CustomerName,
			&item.Email,
			&item.LastSubscription,
			&churnDate,
			&item.LifetimeValue,
		); scanErr != nil {
			return nil, scanErr
		}
		item.ChurnDate = churnDate.Format("2006-01-02")
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []churnedCustomerItem{}
	}
	return items, nil
}

func buildMonthBuckets(startDate, endDate time.Time) []time.Time {
	cursor := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	last := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())
	months := make([]time.Time, 0)
	for !cursor.After(last) {
		months = append(months, cursor)
		cursor = cursor.AddDate(0, 1, 0)
	}
	if months == nil {
		return []time.Time{}
	}
	return months
}

func monthKey(ts time.Time) string {
	return fmt.Sprintf("%04d-%02d", ts.Year(), ts.Month())
}

func calculateRate(numerator, denominator int) float64 {
	if denominator <= 0 {
		return 0
	}
	val := (float64(numerator) / float64(denominator)) * 100
	return round2(val)
}

func compareRate(current, previous int) float64 {
	if previous <= 0 {
		if current <= 0 {
			return 0
		}
		return 100
	}
	return round2(((float64(current) - float64(previous)) / float64(previous)) * 100)
}

func compareFloatRate(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return round2(((current - previous) / previous) * 100)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func calcTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func parseIntDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func handleServerError(c *gin.Context, err error) {
	setAuditError(c, err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Terjadi kesalahan internal server",
	})
}

func ensureManagementSchema(ctx context.Context) error {
	query := `
		CREATE SCHEMA IF NOT EXISTS management;
		CREATE TABLE IF NOT EXISTS management.audit_logs (
			id BIGSERIAL PRIMARY KEY,
			endpoint TEXT NOT NULL,
			method VARCHAR(16) NOT NULL,
			role VARCHAR(64),
			status_code INT NOT NULL,
			message TEXT,
			payload JSONB,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_management_audit_logs_created_at ON management.audit_logs(created_at DESC);
	`
	_, err := db.Exec(ctx, query)
	return err
}

func isDashboardSourceReady(ctx context.Context) (bool, error) {
	query := `
		SELECT
			to_regclass('subscription.orders') IS NOT NULL
			AND to_regclass('subscription.packages') IS NOT NULL
			AND to_regclass('public.users') IS NOT NULL
	`
	var ready bool
	err := db.QueryRow(ctx, query).Scan(&ready)
	if err != nil {
		return false, err
	}
	return ready, nil
}

func buildEmptyDashboardPayload(periodType string, startDate, endDate time.Time, page, limit int) dashboardCustomerPayload {
	payload := dashboardCustomerPayload{}
	payload.Period.Type = periodType
	payload.Period.StartDate = startDate.Format("2006-01-02")
	payload.Period.EndDate = endDate.Format("2006-01-02")

	payload.Summary.TotalActiveCustomers = 0
	payload.Summary.TotalLoyalCustomers = 0
	payload.Summary.TotalChurnedCustomers = 0
	payload.Summary.ChurnRate = 0
	payload.Summary.ActiveCustomersChangePct = 0
	payload.Summary.ChurnRateChangePct = 0
	payload.Summary.ActiveCustomersLastPeriod = 0
	payload.Summary.ChurnRateLastPeriod = 0

	payload.Charts.MonthlyNewCustomers = []monthlyPoint{}
	payload.Charts.MonthlyChurnCustomers = []monthlyPoint{}
	payload.Charts.ChurnRateTrend = []churnRatePoint{}

	payload.TopLoyalCustomers.Items = []loyalCustomerItem{}
	payload.TopLoyalCustomers.Pagination.Page = page
	payload.TopLoyalCustomers.Pagination.Limit = limit
	payload.TopLoyalCustomers.Pagination.Total = 0
	payload.TopLoyalCustomers.Pagination.TotalPages = 0

	payload.RecentlyChurnedCustomers = []churnedCustomerItem{}
	return payload
}

func buildDummyDashboardPayload(periodType string, startDate, endDate time.Time, page, limit int, search string) dashboardCustomerPayload {
	customers := getDummyCustomers(startDate.Location())
	payload := dashboardCustomerPayload{}
	payload.Period.Type = periodType
	payload.Period.StartDate = startDate.Format("2006-01-02")
	payload.Period.EndDate = endDate.Format("2006-01-02")

	totalActive := 0
	totalLoyal := 0
	totalChurned := 0

	for _, customer := range customers {
		if isCustomerActiveInRange(customer, startDate, endDate) {
			totalActive++
		}

		duration := dummyDurationMonths(customer, endDate)
		if duration >= 12 && isCustomerActiveInRange(customer, startDate, endDate) {
			totalLoyal++
		}

		if isCustomerChurnedInRange(customer, startDate, endDate) {
			totalChurned++
		}
	}

	activeAtStart := countDummyActiveAtDate(customers, startDate)
	churnRate := calculateRate(totalChurned, activeAtStart)

	prevStart, prevEnd := previousPeriod(startDate, endDate)
	prevActive := countDummyActiveInRange(customers, prevStart, prevEnd)
	prevChurned := countDummyChurnedInRange(customers, prevStart, prevEnd)
	prevActiveAtStart := countDummyActiveAtDate(customers, prevStart)
	prevChurnRate := calculateRate(prevChurned, prevActiveAtStart)

	payload.Summary.TotalActiveCustomers = totalActive
	payload.Summary.TotalLoyalCustomers = totalLoyal
	payload.Summary.TotalChurnedCustomers = totalChurned
	payload.Summary.ChurnRate = churnRate
	payload.Summary.ActiveCustomersLastPeriod = prevActive
	payload.Summary.ChurnRateLastPeriod = prevChurnRate
	payload.Summary.ActiveCustomersChangePct = compareRate(totalActive, prevActive)
	payload.Summary.ChurnRateChangePct = compareFloatRate(churnRate, prevChurnRate)

	months := buildMonthBuckets(startDate, endDate)
	payload.Charts.MonthlyNewCustomers = make([]monthlyPoint, 0, len(months))
	payload.Charts.MonthlyChurnCustomers = make([]monthlyPoint, 0, len(months))
	payload.Charts.ChurnRateTrend = make([]churnRatePoint, 0, len(months))

	for _, monthStart := range months {
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
		newTotal := 0
		churnTotal := 0
		for _, customer := range customers {
			if isWithinRange(customer.JoinDate, monthStart, monthEnd) {
				newTotal++
			}
			if isCustomerChurnedInRange(customer, monthStart, monthEnd) {
				churnTotal++
			}
		}

		label := monthStart.Format("Jan 2006")
		activeAtMonthStart := countDummyActiveAtDate(customers, monthStart)

		payload.Charts.MonthlyNewCustomers = append(payload.Charts.MonthlyNewCustomers, monthlyPoint{Month: label, Total: newTotal})
		payload.Charts.MonthlyChurnCustomers = append(payload.Charts.MonthlyChurnCustomers, monthlyPoint{Month: label, Total: churnTotal})
		payload.Charts.ChurnRateTrend = append(payload.Charts.ChurnRateTrend, churnRatePoint{Month: label, Rate: calculateRate(churnTotal, activeAtMonthStart)})
	}

	searchLower := strings.ToLower(strings.TrimSpace(search))
	loyalItems := make([]loyalCustomerItem, 0)
	for _, customer := range customers {
		if !isCustomerActiveInRange(customer, startDate, endDate) && !isCustomerChurnedInRange(customer, startDate, endDate) {
			continue
		}

		duration := dummyDurationMonths(customer, endDate)
		if duration < 12 {
			continue
		}

		if searchLower != "" {
			name := strings.ToLower(customer.CustomerName)
			email := strings.ToLower(customer.Email)
			if !strings.Contains(name, searchLower) && !strings.Contains(email, searchLower) {
				continue
			}
		}

		loyalItems = append(loyalItems, loyalCustomerItem{
			CustomerID:   customer.CustomerID,
			CustomerName: customer.CustomerName,
			Email:        customer.Email,
			Duration:     duration,
			Transactions: customer.Transactions,
			TotalSpent:   customer.TotalSpent,
			LastActive:   customer.LastActive.Format(time.RFC3339),
		})
	}

	sort.Slice(loyalItems, func(i, j int) bool {
		if loyalItems[i].Duration == loyalItems[j].Duration {
			return loyalItems[i].TotalSpent > loyalItems[j].TotalSpent
		}
		return loyalItems[i].Duration > loyalItems[j].Duration
	})

	totalLoyalItems := len(loyalItems)
	startIdx := (page - 1) * limit
	if startIdx > totalLoyalItems {
		startIdx = totalLoyalItems
	}
	endIdx := startIdx + limit
	if endIdx > totalLoyalItems {
		endIdx = totalLoyalItems
	}

	payload.TopLoyalCustomers.Items = loyalItems[startIdx:endIdx]
	payload.TopLoyalCustomers.Pagination.Page = page
	payload.TopLoyalCustomers.Pagination.Limit = limit
	payload.TopLoyalCustomers.Pagination.Total = totalLoyalItems
	payload.TopLoyalCustomers.Pagination.TotalPages = calcTotalPages(totalLoyalItems, limit)

	recentlyChurned := make([]churnedCustomerItem, 0)
	for _, customer := range customers {
		if !isCustomerChurnedInRange(customer, startDate, endDate) {
			continue
		}
		recentlyChurned = append(recentlyChurned, churnedCustomerItem{
			CustomerID:       customer.CustomerID,
			CustomerName:     customer.CustomerName,
			Email:            customer.Email,
			LastSubscription: customer.LastSubscription,
			ChurnDate:        customer.ChurnDate.Format("2006-01-02"),
			LifetimeValue:    customer.TotalSpent,
		})
	}

	sort.Slice(recentlyChurned, func(i, j int) bool {
		return recentlyChurned[i].ChurnDate > recentlyChurned[j].ChurnDate
	})
	if len(recentlyChurned) > 10 {
		recentlyChurned = recentlyChurned[:10]
	}
	payload.RecentlyChurnedCustomers = recentlyChurned

	if payload.TopLoyalCustomers.Items == nil {
		payload.TopLoyalCustomers.Items = []loyalCustomerItem{}
	}
	if payload.RecentlyChurnedCustomers == nil {
		payload.RecentlyChurnedCustomers = []churnedCustomerItem{}
	}

	return payload
}

func buildDummyCustomerDetail(id string, now time.Time) (customerDetailPayload, bool) {
	for _, customer := range getDummyCustomers(now.Location()) {
		if customer.CustomerID != id {
			continue
		}

		detail := customerDetailPayload{
			CustomerID:   customer.CustomerID,
			CustomerName: customer.CustomerName,
			Email:        customer.Email,
			Duration:     dummyDurationMonths(customer, now),
			Transactions: customer.Transactions,
			TotalSpent:   customer.TotalSpent,
			LastActive:   customer.LastActive.Format(time.RFC3339),
		}
		return detail, true
	}

	return customerDetailPayload{}, false
}

func getDummyCustomers(loc *time.Location) []dummyCustomer {
	if loc == nil {
		loc = time.Local
	}

	date := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 10, 0, 0, 0, loc)
	}
	pDate := func(y int, m time.Month, d int) *time.Time {
		v := date(y, m, d)
		return &v
	}

	return []dummyCustomer{
		{CustomerID: "CUS-001", CustomerName: "Andi Pratama", Email: "andi.pratama@thinktala.com", JoinDate: date(2024, time.January, 10), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 12), Transactions: 22, TotalSpent: 32500000, LastSubscription: "Elite Plan"},
		{CustomerID: "CUS-002", CustomerName: "Bunga Lestari", Email: "bunga.lestari@thinktala.com", JoinDate: date(2024, time.May, 20), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 11), Transactions: 19, TotalSpent: 27400000, LastSubscription: "Pro Plan"},
		{CustomerID: "CUS-003", CustomerName: "Chandra Wijaya", Email: "chandra.wijaya@thinktala.com", JoinDate: date(2025, time.February, 11), ActiveUntil: date(2026, time.November, 30), LastActive: date(2026, time.April, 10), Transactions: 14, TotalSpent: 18600000, LastSubscription: "Pro Plan"},
		{CustomerID: "CUS-004", CustomerName: "Dinda Maharani", Email: "dinda.maharani@thinktala.com", JoinDate: date(2025, time.August, 8), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 8), Transactions: 11, TotalSpent: 14200000, LastSubscription: "Pro Plan"},
		{CustomerID: "CUS-005", CustomerName: "Eka Firmansyah", Email: "eka.firmansyah@thinktala.com", JoinDate: date(2025, time.November, 25), ActiveUntil: date(2026, time.March, 18), LastActive: date(2026, time.March, 18), Transactions: 6, TotalSpent: 7200000, LastSubscription: "Pro Plan", ChurnDate: pDate(2026, time.March, 18)},
		{CustomerID: "CUS-006", CustomerName: "Farhan Aditya", Email: "farhan.aditya@thinktala.com", JoinDate: date(2026, time.January, 6), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 13), Transactions: 5, TotalSpent: 4900000, LastSubscription: "Pro Plan"},
		{CustomerID: "CUS-007", CustomerName: "Gita Novita", Email: "gita.novita@thinktala.com", JoinDate: date(2026, time.February, 17), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 9), Transactions: 4, TotalSpent: 4400000, LastSubscription: "Pro Plan"},
		{CustomerID: "CUS-008", CustomerName: "Hendra Saputra", Email: "hendra.saputra@thinktala.com", JoinDate: date(2026, time.March, 3), ActiveUntil: date(2026, time.April, 12), LastActive: date(2026, time.April, 12), Transactions: 3, TotalSpent: 3600000, LastSubscription: "Free Plan", ChurnDate: pDate(2026, time.April, 12)},
		{CustomerID: "CUS-009", CustomerName: "Intan Permata", Email: "intan.permata@thinktala.com", JoinDate: date(2026, time.April, 4), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 14), Transactions: 2, TotalSpent: 2800000, LastSubscription: "Pro Plan"},
		{CustomerID: "CUS-010", CustomerName: "Joko Santoso", Email: "joko.santoso@thinktala.com", JoinDate: date(2026, time.April, 19), ActiveUntil: date(2027, time.December, 31), LastActive: date(2026, time.April, 20), Transactions: 1, TotalSpent: 1500000, LastSubscription: "Free Plan"},
		{CustomerID: "CUS-011", CustomerName: "Karin Amelia", Email: "karin.amelia@thinktala.com", JoinDate: date(2025, time.October, 7), ActiveUntil: date(2026, time.February, 28), LastActive: date(2026, time.February, 28), Transactions: 4, TotalSpent: 5200000, LastSubscription: "Pro Plan", ChurnDate: pDate(2026, time.February, 28)},
		{CustomerID: "CUS-012", CustomerName: "Lala Puspita", Email: "lala.puspita@thinktala.com", JoinDate: date(2024, time.September, 14), ActiveUntil: date(2026, time.January, 30), LastActive: date(2026, time.January, 30), Transactions: 13, TotalSpent: 16800000, LastSubscription: "Elite Plan", ChurnDate: pDate(2026, time.January, 30)},
	}
}

func isWithinRange(ts, startDate, endDate time.Time) bool {
	return (ts.Equal(startDate) || ts.After(startDate)) && (ts.Equal(endDate) || ts.Before(endDate))
}

func isCustomerActiveInRange(customer dummyCustomer, startDate, endDate time.Time) bool {
	return (customer.JoinDate.Equal(endDate) || customer.JoinDate.Before(endDate)) &&
		(customer.ActiveUntil.Equal(startDate) || customer.ActiveUntil.After(startDate))
}

func isCustomerChurnedInRange(customer dummyCustomer, startDate, endDate time.Time) bool {
	if customer.ChurnDate == nil {
		return false
	}
	return isWithinRange(*customer.ChurnDate, startDate, endDate)
}

func dummyDurationMonths(customer dummyCustomer, endDate time.Time) int {
	effectiveEnd := endDate
	if customer.ActiveUntil.Before(effectiveEnd) {
		effectiveEnd = customer.ActiveUntil
	}
	if effectiveEnd.Before(customer.JoinDate) {
		return 0
	}

	months := (effectiveEnd.Year()-customer.JoinDate.Year())*12 + int(effectiveEnd.Month()-customer.JoinDate.Month())
	if effectiveEnd.Day() >= customer.JoinDate.Day() {
		months++
	}
	if months < 1 {
		months = 1
	}
	return months
}

func countDummyActiveAtDate(customers []dummyCustomer, at time.Time) int {
	total := 0
	for _, customer := range customers {
		if (customer.JoinDate.Equal(at) || customer.JoinDate.Before(at)) &&
			(customer.ActiveUntil.Equal(at) || customer.ActiveUntil.After(at)) {
			total++
		}
	}
	return total
}

func countDummyActiveInRange(customers []dummyCustomer, startDate, endDate time.Time) int {
	total := 0
	for _, customer := range customers {
		if isCustomerActiveInRange(customer, startDate, endDate) {
			total++
		}
	}
	return total
}

func countDummyChurnedInRange(customers []dummyCustomer, startDate, endDate time.Time) int {
	total := 0
	for _, customer := range customers {
		if isCustomerChurnedInRange(customer, startDate, endDate) {
			total++
		}
	}
	return total
}

func auditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		payload := requestAudit{
			Endpoint: c.Request.URL.Path,
			Method:   c.Request.Method,
			Role:     strings.ToUpper(strings.TrimSpace(c.GetHeader("X-User-Role"))),
			Status:   c.Writer.Status(),
			Message:  strings.TrimSpace(c.GetString("audit_error")),
			Payload:  c.Request.URL.RawQuery,
		}
		if err := saveAuditLog(c.Request.Context(), payload); err != nil {
			log.Printf("[MANAGEMENT] gagal menyimpan audit log: %v", err)
		}
	}
}

func setAuditError(c *gin.Context, msg string) {
	if msg == "" {
		return
	}
	c.Set("audit_error", msg)
}

func saveAuditLog(ctx context.Context, audit requestAudit) error {
	query := `
		INSERT INTO management.audit_logs (endpoint, method, role, status_code, message, payload)
		VALUES ($1, $2, $3, $4, $5, CASE WHEN $6 = '' THEN NULL ELSE to_jsonb($6::text) END)
	`
	_, err := db.Exec(ctx, query, audit.Endpoint, audit.Method, audit.Role, audit.Status, audit.Message, audit.Payload)
	return err
}

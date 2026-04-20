package dashboard

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	repo              Repository
	graceDays         int
	loyalMinimumMonth int
	loyalMinimumTx    int
	churnListLimit    int
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:              repo,
		graceDays:         30,
		loyalMinimumMonth: 3,
		loyalMinimumTx:    2,
		churnListLimit:    10,
	}
}

func (s *Service) GetDashboardCustomers(ctx context.Context, q dashboardQuery) (dashboardCustomerPayload, error) {
	payload := buildEmptyCustomerDashboardPayload(q.PeriodType, q.StartDate, q.EndDate, q.Page, q.Limit)

	ready, err := s.repo.IsCustomerSourceReady(ctx)
	if err != nil {
		return payload, err
	}
	if !ready {
		return payload, nil
	}

	totalActive, err := s.repo.CountActiveCustomers(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return payload, err
	}

	totalLoyal, err := s.repo.CountLoyalCustomers(ctx, q.StartDate, q.EndDate, s.loyalMinimumMonth, s.loyalMinimumTx)
	if err != nil {
		return payload, err
	}

	totalChurned, err := s.repo.CountChurnedCustomers(ctx, q.StartDate, q.EndDate, s.graceDays)
	if err != nil {
		return payload, err
	}

	activeAtStart, err := s.repo.CountActiveAtDate(ctx, q.StartDate)
	if err != nil {
		return payload, err
	}
	churnRate := calculateRate(totalChurned, activeAtStart)

	prevStart, prevEnd := previousPeriod(q.StartDate, q.EndDate)
	prevActive, err := s.repo.CountActiveCustomers(ctx, prevStart, prevEnd)
	if err != nil {
		return payload, err
	}
	prevChurned, err := s.repo.CountChurnedCustomers(ctx, prevStart, prevEnd, s.graceDays)
	if err != nil {
		return payload, err
	}
	prevActiveAtStart, err := s.repo.CountActiveAtDate(ctx, prevStart)
	if err != nil {
		return payload, err
	}
	prevChurnRate := calculateRate(prevChurned, prevActiveAtStart)

	payload.Summary.TotalActiveCustomers = totalActive
	payload.Summary.TotalLoyalCustomers = totalLoyal
	payload.Summary.TotalChurnedCustomers = totalChurned
	payload.Summary.ChurnRate = churnRate
	payload.Summary.ActiveCustomersLastPeriod = prevActive
	payload.Summary.ChurnRateLastPeriod = prevChurnRate
	payload.Summary.ActiveCustomersChangePct = compareRate(totalActive, prevActive)
	payload.Summary.ChurnRateChangePct = compareFloatRate(churnRate, prevChurnRate)

	monthlyNewMap, err := s.repo.FetchMonthlyNewCustomers(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return payload, err
	}
	monthlyChurnMap, err := s.repo.FetchMonthlyChurnCustomers(ctx, q.StartDate, q.EndDate, s.graceDays)
	if err != nil {
		return payload, err
	}

	months := buildMonthBuckets(q.StartDate, q.EndDate)
	payload.Charts.MonthlyNewCustomers = make([]monthlyPoint, 0, len(months))
	payload.Charts.MonthlyChurnCustomers = make([]monthlyPoint, 0, len(months))
	payload.Charts.ChurnRateTrend = make([]churnRatePoint, 0, len(months))

	for _, monthStart := range months {
		key := monthKey(monthStart)
		label := monthStart.Format("Jan 2006")
		newTotal := monthlyNewMap[key]
		churnTotal := monthlyChurnMap[key]

		activeAtMonthStart, err := s.repo.CountActiveAtDate(ctx, monthStart)
		if err != nil {
			return payload, err
		}

		payload.Charts.MonthlyNewCustomers = append(payload.Charts.MonthlyNewCustomers, monthlyPoint{
			Month: label,
			Total: newTotal,
		})
		payload.Charts.MonthlyChurnCustomers = append(payload.Charts.MonthlyChurnCustomers, monthlyPoint{
			Month: label,
			Total: churnTotal,
		})
		payload.Charts.ChurnRateTrend = append(payload.Charts.ChurnRateTrend, churnRatePoint{
			Month: label,
			Rate:  calculateRate(churnTotal, activeAtMonthStart),
		})
	}

	loyalItems, loyalTotal, err := s.repo.FetchTopLoyalCustomers(ctx, q.StartDate, q.EndDate, s.loyalMinimumMonth, s.loyalMinimumTx, q.Search, q.Page, q.Limit)
	if err != nil {
		return payload, err
	}
	payload.TopLoyalCustomers.Items = loyalItems
	payload.TopLoyalCustomers.Pagination.Page = q.Page
	payload.TopLoyalCustomers.Pagination.Limit = q.Limit
	payload.TopLoyalCustomers.Pagination.Total = loyalTotal
	payload.TopLoyalCustomers.Pagination.TotalPages = calcTotalPages(loyalTotal, q.Limit)

	recentlyChurned, err := s.repo.FetchRecentlyChurnedCustomers(ctx, q.StartDate, q.EndDate, s.graceDays, s.churnListLimit)
	if err != nil {
		return payload, err
	}
	payload.RecentlyChurnedCustomers = recentlyChurned

	if payload.TopLoyalCustomers.Items == nil {
		payload.TopLoyalCustomers.Items = []loyalCustomerItem{}
	}
	if payload.RecentlyChurnedCustomers == nil {
		payload.RecentlyChurnedCustomers = []churnedCustomerItem{}
	}

	return payload, nil
}

func (s *Service) GetCustomerDetail(ctx context.Context, customerID string) (customerDetailPayload, bool, error) {
	ready, err := s.repo.IsCustomerSourceReady(ctx)
	if err != nil {
		return customerDetailPayload{}, false, err
	}
	if !ready {
		return customerDetailPayload{}, false, nil
	}

	return s.repo.GetCustomerDetail(ctx, customerID, time.Now())
}

func (s *Service) GetDashboardPackages(ctx context.Context, q dashboardQuery) (dashboardPackagePayload, error) {
	payload := buildEmptyPackageDashboardPayload(q.PeriodType, q.StartDate, q.EndDate, q.Page, q.Limit)

	ready, err := s.repo.IsPackageSourceReady(ctx)
	if err != nil {
		return payload, err
	}
	if !ready {
		return payload, nil
	}

	catalog, err := s.repo.ListPackages(ctx)
	if err != nil {
		return payload, err
	}

	currentMetrics, err := s.repo.AggregatePackageMetrics(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return payload, err
	}

	prevStart, prevEnd := previousPeriod(q.StartDate, q.EndDate)
	prevMetrics, err := s.repo.AggregatePackageMetrics(ctx, prevStart, prevEnd)
	if err != nil {
		return payload, err
	}

	for id, metric := range currentMetrics {
		exists := false
		for _, item := range catalog {
			if item.ID == id {
				exists = true
				break
			}
		}
		if !exists {
			catalog = append(catalog, packageCatalogItem{ID: id, Name: metric.PackageName})
		}
	}

	totalRevenueAll := 0.0
	for _, metric := range currentMetrics {
		totalRevenueAll += metric.TotalRevenue
	}

	metricsList := make([]packagePerformanceItem, 0, len(catalog))

	for _, item := range catalog {
		current := currentMetrics[item.ID]
		prev := prevMetrics[item.ID]

		currTx := current.TotalTransactions
		currRevenue := current.TotalRevenue
		if current.PackageName == "" {
			current.PackageName = item.Name
		}

		prevTx := prev.TotalTransactions
		growth := compareRate(currTx, prevTx)
		contribution := 0.0
		if totalRevenueAll > 0 {
			contribution = round2((currRevenue / totalRevenueAll) * 100)
		}

		if currTx <= 0 {
			continue
		}

		row := packagePerformanceItem{
			PackageID:              item.ID,
			PackageName:            current.PackageName,
			TotalTransactions:      currTx,
			TotalRevenue:           currRevenue,
			PercentageContribution: contribution,
			GrowthRate:             growth,
		}
		metricsList = append(metricsList, row)
	}

	if len(metricsList) > 0 {
		mostSoldList := append([]packagePerformanceItem{}, metricsList...)
		sort.Slice(mostSoldList, func(i, j int) bool {
			if mostSoldList[i].TotalTransactions == mostSoldList[j].TotalTransactions {
				if mostSoldList[i].TotalRevenue == mostSoldList[j].TotalRevenue {
					return mostSoldList[i].PackageName < mostSoldList[j].PackageName
				}
				return mostSoldList[i].TotalRevenue > mostSoldList[j].TotalRevenue
			}
			return mostSoldList[i].TotalTransactions > mostSoldList[j].TotalTransactions
		})
		topSales := mostSoldList[0].TotalTransactions
		tiedSalesCount := 0
		for _, item := range mostSoldList {
			if item.TotalTransactions == topSales {
				tiedSalesCount++
			}
		}
		payload.Summary.MostSoldPackage.PackageID = mostSoldList[0].PackageID
		payload.Summary.MostSoldPackage.TotalSales = topSales
		if tiedSalesCount > 1 {
			payload.Summary.MostSoldPackage.PackageName = "Tie (" + strconv.Itoa(tiedSalesCount) + " Packages)"
		} else {
			payload.Summary.MostSoldPackage.PackageName = mostSoldList[0].PackageName
		}

		highestRevenueList := append([]packagePerformanceItem{}, metricsList...)
		sort.Slice(highestRevenueList, func(i, j int) bool {
			if highestRevenueList[i].TotalRevenue == highestRevenueList[j].TotalRevenue {
				if highestRevenueList[i].TotalTransactions == highestRevenueList[j].TotalTransactions {
					return highestRevenueList[i].PackageName < highestRevenueList[j].PackageName
				}
				return highestRevenueList[i].TotalTransactions > highestRevenueList[j].TotalTransactions
			}
			return highestRevenueList[i].TotalRevenue > highestRevenueList[j].TotalRevenue
		})
		topRevenue := highestRevenueList[0].TotalRevenue
		tiedRevenueCount := 0
		for _, item := range highestRevenueList {
			if item.TotalRevenue == topRevenue {
				tiedRevenueCount++
			}
		}
		payload.Summary.HighestRevenuePackage.PackageID = highestRevenueList[0].PackageID
		payload.Summary.HighestRevenuePackage.TotalRevenue = topRevenue
		if tiedRevenueCount > 1 {
			payload.Summary.HighestRevenuePackage.PackageName = "Tie (" + strconv.Itoa(tiedRevenueCount) + " Packages)"
		} else {
			payload.Summary.HighestRevenuePackage.PackageName = highestRevenueList[0].PackageName
		}

		fastestGrowthList := append([]packagePerformanceItem{}, metricsList...)
		sort.Slice(fastestGrowthList, func(i, j int) bool {
			if fastestGrowthList[i].GrowthRate == fastestGrowthList[j].GrowthRate {
				if fastestGrowthList[i].TotalTransactions == fastestGrowthList[j].TotalTransactions {
					if fastestGrowthList[i].TotalRevenue == fastestGrowthList[j].TotalRevenue {
						return fastestGrowthList[i].PackageName < fastestGrowthList[j].PackageName
					}
					return fastestGrowthList[i].TotalRevenue > fastestGrowthList[j].TotalRevenue
				}
				return fastestGrowthList[i].TotalTransactions > fastestGrowthList[j].TotalTransactions
			}
			return fastestGrowthList[i].GrowthRate > fastestGrowthList[j].GrowthRate
		})
		topGrowth := fastestGrowthList[0].GrowthRate
		tiedGrowthCount := 0
		for _, item := range fastestGrowthList {
			if item.GrowthRate == topGrowth {
				tiedGrowthCount++
			}
		}
		payload.Summary.FastestGrowthPackage.PackageID = fastestGrowthList[0].PackageID
		payload.Summary.FastestGrowthPackage.GrowthRate = topGrowth
		if tiedGrowthCount > 1 {
			payload.Summary.FastestGrowthPackage.PackageName = "Tie (" + strconv.Itoa(tiedGrowthCount) + " Packages)"
		} else {
			payload.Summary.FastestGrowthPackage.PackageName = fastestGrowthList[0].PackageName
		}
	} else {
		payload.Summary.MostSoldPackage.PackageName = "-"
		payload.Summary.HighestRevenuePackage.PackageName = "-"
		payload.Summary.FastestGrowthPackage.PackageName = "-"
		payload.Summary.FastestGrowthPackage.GrowthRate = 0
	}

	payload.Summary.TotalRevenueAllPackage = round2(totalRevenueAll)

	searchLower := strings.ToLower(strings.TrimSpace(q.Search))
	filtered := make([]packagePerformanceItem, 0, len(metricsList))
	for _, item := range metricsList {
		if searchLower != "" && !strings.Contains(strings.ToLower(item.PackageName), searchLower) {
			continue
		}
		filtered = append(filtered, item)
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].TotalTransactions == filtered[j].TotalTransactions {
			return filtered[i].TotalRevenue > filtered[j].TotalRevenue
		}
		return filtered[i].TotalTransactions > filtered[j].TotalTransactions
	})

	totalItems := len(filtered)
	startIdx := (q.Page - 1) * q.Limit
	if startIdx > totalItems {
		startIdx = totalItems
	}
	endIdx := startIdx + q.Limit
	if endIdx > totalItems {
		endIdx = totalItems
	}

	payload.Packages.Items = filtered[startIdx:endIdx]
	payload.Packages.Pagination.Page = q.Page
	payload.Packages.Pagination.Limit = q.Limit
	payload.Packages.Pagination.Total = totalItems
	payload.Packages.Pagination.TotalPages = calcTotalPages(totalItems, q.Limit)

	distribution := make([]packageSalesDistributionPoint, 0, len(metricsList))
	contribution := make([]packageRevenueContributionPoint, 0, len(metricsList))
	for _, item := range metricsList {
		distribution = append(distribution, packageSalesDistributionPoint{
			PackageID:   item.PackageID,
			PackageName: item.PackageName,
			TotalSales:  item.TotalTransactions,
		})
		contribution = append(contribution, packageRevenueContributionPoint{
			PackageID:              item.PackageID,
			PackageName:            item.PackageName,
			Revenue:                item.TotalRevenue,
			PercentageContribution: item.PercentageContribution,
		})
	}
	sort.Slice(distribution, func(i, j int) bool {
		if distribution[i].TotalSales == distribution[j].TotalSales {
			return distribution[i].PackageName < distribution[j].PackageName
		}
		return distribution[i].TotalSales > distribution[j].TotalSales
	})
	sort.Slice(contribution, func(i, j int) bool {
		if contribution[i].Revenue == contribution[j].Revenue {
			return contribution[i].PackageName < contribution[j].PackageName
		}
		return contribution[i].Revenue > contribution[j].Revenue
	})
	payload.Charts.PackageSalesDistribution = distribution
	payload.Charts.RevenueContribution = contribution

	monthlyMetrics, err := s.repo.AggregatePackageMonthlyMetrics(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return payload, err
	}

	bucketMap := map[string]packageSalesTrendPoint{}
	for _, month := range buildMonthBuckets(q.StartDate, q.EndDate) {
		key := monthKey(month)
		bucketMap[key] = packageSalesTrendPoint{Period: month.Format("Jan 2006")}
	}

	for _, row := range monthlyMetrics {
		key := monthKey(row.MonthStart)
		point, ok := bucketMap[key]
		if !ok {
			point = packageSalesTrendPoint{Period: row.MonthStart.Format("Jan 2006")}
		}

		switch classifyPackageBucket(row.PackageID, row.PackageName) {
		case "premium":
			point.Premium += row.TotalSales
		case "enterprise":
			point.Enterprise += row.TotalSales
		default:
			point.Starter += row.TotalSales
		}
		bucketMap[key] = point
	}

	trend := make([]packageSalesTrendPoint, 0, len(bucketMap))
	for _, month := range buildMonthBuckets(q.StartDate, q.EndDate) {
		trend = append(trend, bucketMap[monthKey(month)])
	}
	payload.Charts.SalesTrend = trend

	if payload.Packages.Items == nil {
		payload.Packages.Items = []packagePerformanceItem{}
	}
	if payload.Charts.PackageSalesDistribution == nil {
		payload.Charts.PackageSalesDistribution = []packageSalesDistributionPoint{}
	}
	if payload.Charts.RevenueContribution == nil {
		payload.Charts.RevenueContribution = []packageRevenueContributionPoint{}
	}
	if payload.Charts.SalesTrend == nil {
		payload.Charts.SalesTrend = []packageSalesTrendPoint{}
	}

	return payload, nil
}

func (s *Service) GetPackageDetail(ctx context.Context, packageID string, q dashboardQuery) (packageDetailPayload, bool, error) {
	ready, err := s.repo.IsPackageSourceReady(ctx)
	if err != nil {
		return packageDetailPayload{}, false, err
	}
	if !ready {
		return packageDetailPayload{}, false, nil
	}

	selected, err := s.repo.GetPackageByID(ctx, packageID)
	if err != nil {
		return packageDetailPayload{}, false, err
	}
	if selected == nil {
		return packageDetailPayload{}, false, nil
	}

	currentMetrics, err := s.repo.AggregatePackageMetrics(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return packageDetailPayload{}, false, err
	}
	prevStart, prevEnd := previousPeriod(q.StartDate, q.EndDate)
	prevMetrics, err := s.repo.AggregatePackageMetrics(ctx, prevStart, prevEnd)
	if err != nil {
		return packageDetailPayload{}, false, err
	}

	totalRevenueAll, err := s.repo.GetPackageTotalRevenue(ctx, q.StartDate, q.EndDate)
	if err != nil {
		return packageDetailPayload{}, false, err
	}

	current := currentMetrics[selected.ID]
	prev := prevMetrics[selected.ID]

	totalTransactions := current.TotalTransactions
	totalRevenue := current.TotalRevenue
	if current.PackageName == "" {
		current.PackageName = selected.Name
	}

	growth := compareRate(totalTransactions, prev.TotalTransactions)
	contribution := 0.0
	if totalRevenueAll > 0 {
		contribution = round2((totalRevenue / totalRevenueAll) * 100)
	}

	monthlyMap, err := s.repo.AggregatePackageMonthlyMetricsByPackage(ctx, selected.ID, q.StartDate, q.EndDate)
	if err != nil {
		return packageDetailPayload{}, false, err
	}

	months := buildMonthBuckets(q.StartDate, q.EndDate)
	trend := make([]packageDetailTrendPoint, 0, len(months))
	for _, monthStart := range months {
		key := monthKey(monthStart)
		trend = append(trend, packageDetailTrendPoint{
			Period: monthStart.Format("Jan 2006"),
			Sales:  monthlyMap[key],
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
			"type":       q.PeriodType,
			"start_date": q.StartDate.Format("2006-01-02"),
			"end_date":   q.EndDate.Format("2006-01-02"),
		},
		Trend: trend,
	}

	if payload.Trend == nil {
		payload.Trend = []packageDetailTrendPoint{}
	}

	return payload, true, nil
}

func buildEmptyCustomerDashboardPayload(periodType string, startDate, endDate time.Time, page, limit int) dashboardCustomerPayload {
	payload := dashboardCustomerPayload{}
	payload.Period.Type = periodType
	payload.Period.StartDate = startDate.Format("2006-01-02")
	payload.Period.EndDate = endDate.Format("2006-01-02")

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

func buildEmptyPackageDashboardPayload(periodType string, startDate, endDate time.Time, page, limit int) dashboardPackagePayload {
	payload := dashboardPackagePayload{}
	payload.Period.Type = periodType
	payload.Period.StartDate = startDate.Format("2006-01-02")
	payload.Period.EndDate = endDate.Format("2006-01-02")
	payload.Summary.MostSoldPackage.PackageName = "-"
	payload.Summary.HighestRevenuePackage.PackageName = "-"
	payload.Summary.FastestGrowthPackage.PackageName = "-"
	payload.Summary.PeriodLabel = packagePeriodLabel(periodType)
	payload.Charts.PackageSalesDistribution = []packageSalesDistributionPoint{}
	payload.Charts.RevenueContribution = []packageRevenueContributionPoint{}
	payload.Charts.SalesTrend = []packageSalesTrendPoint{}
	payload.Packages.Items = []packagePerformanceItem{}
	payload.Packages.Pagination.Page = page
	payload.Packages.Pagination.Limit = limit
	payload.Packages.Pagination.Total = 0
	payload.Packages.Pagination.TotalPages = 0
	return payload
}

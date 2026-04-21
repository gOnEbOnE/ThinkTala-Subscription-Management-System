package dashboard

import "time"

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

type requestAudit struct {
	Endpoint string
	Method   string
	Role     string
	Status   int
	Message  string
	Payload  string
}

type packageCatalogItem struct {
	ID   string
	Name string
}

type packageMetric struct {
	PackageID         string
	PackageName       string
	TotalTransactions int
	TotalRevenue      float64
}

type packageMonthlyMetric struct {
	MonthStart  time.Time
	PackageID   string
	PackageName string
	TotalSales  int
}

type dashboardQuery struct {
	PeriodType string
	StartDate  time.Time
	EndDate    time.Time
	Page       int
	Limit      int
	Search     string
}

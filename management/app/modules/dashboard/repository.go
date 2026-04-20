package dashboard

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	IsCustomerSourceReady(ctx context.Context) (bool, error)
	IsPackageSourceReady(ctx context.Context) (bool, error)
	CountActiveCustomers(ctx context.Context, startDate, endDate time.Time) (int, error)
	CountActiveAtDate(ctx context.Context, at time.Time) (int, error)
	CountChurnedCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (int, error)
	CountLoyalCustomers(ctx context.Context, startDate, endDate time.Time, minimumMonths, minimumTransactions int) (int, error)
	FetchMonthlyNewCustomers(ctx context.Context, startDate, endDate time.Time) (map[string]int, error)
	FetchMonthlyChurnCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (map[string]int, error)
	FetchTopLoyalCustomers(ctx context.Context, startDate, endDate time.Time, minimumMonths, minimumTransactions int, search string, page, limit int) ([]loyalCustomerItem, int, error)
	FetchRecentlyChurnedCustomers(ctx context.Context, startDate, endDate time.Time, graceDays, maxItems int) ([]churnedCustomerItem, error)
	GetCustomerDetail(ctx context.Context, customerID string, endDate time.Time) (customerDetailPayload, bool, error)
	ListPackages(ctx context.Context) ([]packageCatalogItem, error)
	GetPackageByID(ctx context.Context, packageID string) (*packageCatalogItem, error)
	AggregatePackageMetrics(ctx context.Context, startDate, endDate time.Time) (map[string]packageMetric, error)
	AggregatePackageMonthlyMetrics(ctx context.Context, startDate, endDate time.Time) ([]packageMonthlyMetric, error)
	AggregatePackageMonthlyMetricsByPackage(ctx context.Context, packageID string, startDate, endDate time.Time) (map[string]int, error)
	GetPackageTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) IsCustomerSourceReady(ctx context.Context) (bool, error) {
	query := `
		SELECT
			to_regclass('subscription.orders') IS NOT NULL
			AND to_regclass('subscription.packages') IS NOT NULL
	`
	var ready bool
	if err := r.db.QueryRow(ctx, query).Scan(&ready); err != nil {
		return false, err
	}
	return ready, nil
}

func (r *repository) IsPackageSourceReady(ctx context.Context) (bool, error) {
	query := `
		SELECT
			to_regclass('subscription.orders') IS NOT NULL
			AND to_regclass('subscription.packages') IS NOT NULL
	`
	var ready bool
	if err := r.db.QueryRow(ctx, query).Scan(&ready); err != nil {
		return false, err
	}
	return ready, nil
}

func (r *repository) CountActiveCustomers(ctx context.Context, startDate, endDate time.Time) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		)
		SELECT COUNT(DISTINCT user_id)
		FROM paid_orders
		WHERE created_at <= $2
		  AND expiry_at >= $1
	`

	var total int
	err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&total)
	return total, err
}

func (r *repository) CountActiveAtDate(ctx context.Context, at time.Time) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		)
		SELECT COUNT(DISTINCT user_id)
		FROM paid_orders
		WHERE created_at <= $1
		  AND expiry_at >= $1
	`

	var total int
	err := r.db.QueryRow(ctx, query, at).Scan(&total)
	return total, err
}

func (r *repository) CountChurnedCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_last_expiry AS (
			SELECT
				user_id,
				MAX(expiry_at) AS last_expiry
			FROM paid_orders
			GROUP BY user_id
		),
		churn_events AS (
			SELECT ule.user_id
			FROM user_last_expiry ule
			WHERE ule.last_expiry >= $1
			  AND ule.last_expiry <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = ule.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > ule.last_expiry
				  AND r.created_at <= (ule.last_expiry + make_interval(days => $3))
			)
		)
		SELECT COUNT(*) FROM churn_events
	`

	var total int
	err := r.db.QueryRow(ctx, query, startDate, endDate, graceDays).Scan(&total)
	return total, err
}

func (r *repository) CountLoyalCustomers(ctx context.Context, startDate, endDate time.Time, minimumMonths, minimumTransactions int) (int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS total_months,
				COUNT(*)::int AS transactions,
				MAX(po.expiry_at) AS last_expiry
			FROM paid_orders po
			WHERE po.created_at <= $2
			GROUP BY po.user_id
		)
		SELECT COUNT(*)
		FROM user_agg ua
		WHERE (ua.total_months >= $3 OR ua.transactions >= $4)
		  AND ua.last_expiry >= $1
	`

	var total int
	err := r.db.QueryRow(ctx, query, startDate, endDate, minimumMonths, minimumTransactions).Scan(&total)
	return total, err
}

func (r *repository) FetchMonthlyNewCustomers(ctx context.Context, startDate, endDate time.Time) (map[string]int, error) {
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

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var monthStart time.Time
		var total int
		if err := rows.Scan(&monthStart, &total); err != nil {
			return nil, err
		}
		result[monthKey(monthStart)] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) FetchMonthlyChurnCustomers(ctx context.Context, startDate, endDate time.Time, graceDays int) (map[string]int, error) {
	query := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_last_expiry AS (
			SELECT
				user_id,
				MAX(expiry_at) AS last_expiry
			FROM paid_orders
			GROUP BY user_id
		),
		churn_events AS (
			SELECT ule.user_id, ule.last_expiry
			FROM user_last_expiry ule
			WHERE ule.last_expiry >= $1
			  AND ule.last_expiry <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = ule.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > ule.last_expiry
				  AND r.created_at <= (ule.last_expiry + make_interval(days => $3))
			)
		)
		SELECT date_trunc('month', last_expiry) AS month_start, COUNT(DISTINCT user_id)::int AS total
		FROM churn_events
		GROUP BY 1
		ORDER BY 1
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate, graceDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var monthStart time.Time
		var total int
		if err := rows.Scan(&monthStart, &total); err != nil {
			return nil, err
		}
		result[monthKey(monthStart)] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) FetchTopLoyalCustomers(ctx context.Context, startDate, endDate time.Time, minimumMonths, minimumTransactions int, search string, page, limit int) ([]loyalCustomerItem, int, error) {
	offset := (page - 1) * limit
	search = strings.ToLower(strings.TrimSpace(search))
	usersReady, err := r.hasUsersTable(ctx)
	if err != nil {
		return nil, 0, err
	}

	countQueryWithUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				o.total_price,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(po.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(po.created_at) AS last_active,
				MAX(po.expiry_at) AS last_expiry
			FROM paid_orders po
			WHERE po.created_at <= $1
			GROUP BY po.user_id
		)
		SELECT COUNT(*)
		FROM user_agg ua
		LEFT JOIN users u ON u.id = ua.user_id
		WHERE (ua.duration_months >= $2 OR ua.transactions >= $3)
		  AND ua.last_expiry >= $4
		  AND (
			$5 = ''
			OR LOWER(COALESCE(u.name, '')) LIKE '%' || $5 || '%'
			OR LOWER(COALESCE(u.email, '')) LIKE '%' || $5 || '%'
			OR LOWER(ua.user_id::text) LIKE '%' || $5 || '%'
		  )
	`
	countQueryWithoutUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				o.total_price,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(po.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(po.created_at) AS last_active,
				MAX(po.expiry_at) AS last_expiry
			FROM paid_orders po
			WHERE po.created_at <= $1
			GROUP BY po.user_id
		)
		SELECT COUNT(*)
		FROM user_agg ua
		WHERE (ua.duration_months >= $2 OR ua.transactions >= $3)
		  AND ua.last_expiry >= $4
		  AND (
			$5 = ''
			OR LOWER(ua.user_id::text) LIKE '%' || $5 || '%'
		  )
	`

	var total int
	if usersReady {
		err = r.db.QueryRow(ctx, countQueryWithUsers, endDate, minimumMonths, minimumTransactions, startDate, search).Scan(&total)
	} else {
		err = r.db.QueryRow(ctx, countQueryWithoutUsers, endDate, minimumMonths, minimumTransactions, startDate, search).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	queryWithUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				o.total_price,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(po.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(po.created_at) AS last_active,
				MAX(po.expiry_at) AS last_expiry
			FROM paid_orders po
			WHERE po.created_at <= $1
			GROUP BY po.user_id
		)
		SELECT
			ua.user_id::text,
			COALESCE(NULLIF(TRIM(u.name), ''), 'Customer ' || LEFT(ua.user_id::text, 8)) AS customer_name,
			COALESCE(u.email, '') AS email,
			ua.duration_months,
			ua.transactions,
			ua.total_spent,
			ua.last_active
		FROM user_agg ua
		LEFT JOIN users u ON u.id = ua.user_id
		WHERE (ua.duration_months >= $2 OR ua.transactions >= $3)
		  AND ua.last_expiry >= $4
		  AND (
			$5 = ''
			OR LOWER(COALESCE(u.name, '')) LIKE '%' || $5 || '%'
			OR LOWER(COALESCE(u.email, '')) LIKE '%' || $5 || '%'
			OR LOWER(ua.user_id::text) LIKE '%' || $5 || '%'
		  )
		ORDER BY ua.transactions DESC, ua.duration_months DESC, ua.total_spent DESC, ua.last_active DESC
		LIMIT $6 OFFSET $7
	`
	queryWithoutUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				o.total_price,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(po.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(po.created_at) AS last_active,
				MAX(po.expiry_at) AS last_expiry
			FROM paid_orders po
			WHERE po.created_at <= $1
			GROUP BY po.user_id
		)
		SELECT
			ua.user_id::text,
			('Customer ' || LEFT(ua.user_id::text, 8)) AS customer_name,
			'' AS email,
			ua.duration_months,
			ua.transactions,
			ua.total_spent,
			ua.last_active
		FROM user_agg ua
		WHERE (ua.duration_months >= $2 OR ua.transactions >= $3)
		  AND ua.last_expiry >= $4
		  AND (
			$5 = ''
			OR LOWER(ua.user_id::text) LIKE '%' || $5 || '%'
		  )
		ORDER BY ua.transactions DESC, ua.duration_months DESC, ua.total_spent DESC, ua.last_active DESC
		LIMIT $6 OFFSET $7
	`

	var rows pgx.Rows
	if usersReady {
		rows, err = r.db.Query(ctx, queryWithUsers, endDate, minimumMonths, minimumTransactions, startDate, search, limit, offset)
	} else {
		rows, err = r.db.Query(ctx, queryWithoutUsers, endDate, minimumMonths, minimumTransactions, startDate, search, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]loyalCustomerItem, 0)
	for rows.Next() {
		var item loyalCustomerItem
		var lastActive time.Time
		if err := rows.Scan(
			&item.CustomerID,
			&item.CustomerName,
			&item.Email,
			&item.Duration,
			&item.Transactions,
			&item.TotalSpent,
			&lastActive,
		); err != nil {
			return nil, 0, err
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

func (r *repository) FetchRecentlyChurnedCustomers(ctx context.Context, startDate, endDate time.Time, graceDays, maxItems int) ([]churnedCustomerItem, error) {
	usersReady, err := r.hasUsersTable(ctx)
	if err != nil {
		return nil, err
	}

	queryWithUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				o.total_price,
				COALESCE(p.name, 'Unknown Plan') AS package_name,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		),
		user_last_order AS (
			SELECT
				user_id,
				MAX(expiry_at) AS last_expiry,
				(array_agg(package_name ORDER BY expiry_at DESC))[1] AS last_package_name
			FROM paid_orders
			GROUP BY user_id
		),
		churn_events AS (
			SELECT
				ulo.user_id,
				ulo.last_package_name AS package_name,
				ulo.last_expiry
			FROM user_last_order ulo
			WHERE ulo.last_expiry >= $1
			  AND ulo.last_expiry <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = ulo.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > ulo.last_expiry
				  AND r.created_at <= (ulo.last_expiry + make_interval(days => $3))
			)
		),
		ltv AS (
			SELECT user_id, COALESCE(SUM(total_price), 0)::numeric(18,2) AS lifetime_value
			FROM subscription.orders
			WHERE status = 'PAID'
			GROUP BY user_id
		)
		SELECT
			ce.user_id::text,
			COALESCE(NULLIF(TRIM(u.name), ''), 'Customer ' || LEFT(ce.user_id::text, 8)) AS customer_name,
			COALESCE(u.email, '') AS email,
			ce.package_name,
			ce.last_expiry,
			COALESCE(l.lifetime_value, 0)
		FROM churn_events ce
		LEFT JOIN users u ON u.id = ce.user_id
		LEFT JOIN ltv l ON l.user_id = ce.user_id
		ORDER BY ce.last_expiry DESC
		LIMIT $4
	`
	queryWithoutUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				o.total_price,
				COALESCE(p.name, 'Unknown Plan') AS package_name,
				(o.created_at + make_interval(months => GREATEST(COALESCE(o.duration_months, 1), 1))) AS expiry_at
			FROM subscription.orders o
			LEFT JOIN subscription.packages p ON p.id = o.package_id
			WHERE o.status = 'PAID'
		),
		user_last_order AS (
			SELECT
				user_id,
				MAX(expiry_at) AS last_expiry,
				(array_agg(package_name ORDER BY expiry_at DESC))[1] AS last_package_name
			FROM paid_orders
			GROUP BY user_id
		),
		churn_events AS (
			SELECT
				ulo.user_id,
				ulo.last_package_name AS package_name,
				ulo.last_expiry
			FROM user_last_order ulo
			WHERE ulo.last_expiry >= $1
			  AND ulo.last_expiry <= $2
			  AND NOT EXISTS (
				SELECT 1
				FROM subscription.orders r
				WHERE r.user_id = ulo.user_id
				  AND r.status = 'PAID'
				  AND r.created_at > ulo.last_expiry
				  AND r.created_at <= (ulo.last_expiry + make_interval(days => $3))
			)
		),
		ltv AS (
			SELECT user_id, COALESCE(SUM(total_price), 0)::numeric(18,2) AS lifetime_value
			FROM subscription.orders
			WHERE status = 'PAID'
			GROUP BY user_id
		)
		SELECT
			ce.user_id::text,
			('Customer ' || LEFT(ce.user_id::text, 8)) AS customer_name,
			'' AS email,
			ce.package_name,
			ce.last_expiry,
			COALESCE(l.lifetime_value, 0)
		FROM churn_events ce
		LEFT JOIN ltv l ON l.user_id = ce.user_id
		ORDER BY ce.last_expiry DESC
		LIMIT $4
	`

	var rows pgx.Rows
	if usersReady {
		rows, err = r.db.Query(ctx, queryWithUsers, startDate, endDate, graceDays, maxItems)
	} else {
		rows, err = r.db.Query(ctx, queryWithoutUsers, startDate, endDate, graceDays, maxItems)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]churnedCustomerItem, 0)
	for rows.Next() {
		var item churnedCustomerItem
		var churnDate time.Time
		if err := rows.Scan(
			&item.CustomerID,
			&item.CustomerName,
			&item.Email,
			&item.LastSubscription,
			&churnDate,
			&item.LifetimeValue,
		); err != nil {
			return nil, err
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

func (r *repository) GetCustomerDetail(ctx context.Context, customerID string, endDate time.Time) (customerDetailPayload, bool, error) {
	usersReady, err := r.hasUsersTable(ctx)
	if err != nil {
		return customerDetailPayload{}, false, err
	}

	queryWithUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				o.total_price
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(po.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(po.created_at) AS last_active
			FROM paid_orders po
			WHERE po.created_at <= $2
			GROUP BY po.user_id
		)
		SELECT
			ua.user_id::text,
			COALESCE(NULLIF(TRIM(u.name), ''), 'Customer ' || LEFT(ua.user_id::text, 8)) AS customer_name,
			COALESCE(u.email, '') AS email,
			ua.duration_months,
			ua.transactions,
			ua.total_spent,
			ua.last_active
		FROM user_agg ua
		LEFT JOIN users u ON u.id = ua.user_id
		WHERE LOWER(ua.user_id::text) = LOWER($1)
		LIMIT 1
	`
	queryWithoutUsers := `
		WITH paid_orders AS (
			SELECT
				o.user_id,
				o.created_at,
				GREATEST(COALESCE(o.duration_months, 1), 1) AS duration_months,
				o.total_price
			FROM subscription.orders o
			WHERE o.status = 'PAID'
		),
		user_agg AS (
			SELECT
				po.user_id,
				SUM(po.duration_months)::int AS duration_months,
				COUNT(*)::int AS transactions,
				COALESCE(SUM(po.total_price), 0)::numeric(18,2) AS total_spent,
				MAX(po.created_at) AS last_active
			FROM paid_orders po
			WHERE po.created_at <= $2
			GROUP BY po.user_id
		)
		SELECT
			ua.user_id::text,
			('Customer ' || LEFT(ua.user_id::text, 8)) AS customer_name,
			'' AS email,
			ua.duration_months,
			ua.transactions,
			ua.total_spent,
			ua.last_active
		FROM user_agg ua
		WHERE LOWER(ua.user_id::text) = LOWER($1)
		LIMIT 1
	`

	var payload customerDetailPayload
	var lastActive time.Time
	if usersReady {
		err = r.db.QueryRow(ctx, queryWithUsers, customerID, endDate).Scan(
			&payload.CustomerID,
			&payload.CustomerName,
			&payload.Email,
			&payload.Duration,
			&payload.Transactions,
			&payload.TotalSpent,
			&lastActive,
		)
	} else {
		err = r.db.QueryRow(ctx, queryWithoutUsers, customerID, endDate).Scan(
			&payload.CustomerID,
			&payload.CustomerName,
			&payload.Email,
			&payload.Duration,
			&payload.Transactions,
			&payload.TotalSpent,
			&lastActive,
		)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return customerDetailPayload{}, false, nil
		}
		return customerDetailPayload{}, false, err
	}
	payload.LastActive = lastActive.Format(time.RFC3339)
	return payload, true, nil
}

func (r *repository) hasUsersTable(ctx context.Context) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT to_regclass('public.users') IS NOT NULL`).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) ListPackages(ctx context.Context) ([]packageCatalogItem, error) {
	query := `
		SELECT id::text, COALESCE(name, 'Unknown Package')
		FROM subscription.packages
		WHERE status != 'DELETED'
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]packageCatalogItem, 0)
	for rows.Next() {
		var item packageCatalogItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []packageCatalogItem{}
	}
	return items, nil
}

func (r *repository) GetPackageByID(ctx context.Context, packageID string) (*packageCatalogItem, error) {
	query := `
		SELECT id::text, COALESCE(name, 'Unknown Package')
		FROM subscription.packages
		WHERE LOWER(id::text) = LOWER($1)
		  AND status != 'DELETED'
		LIMIT 1
	`

	var item packageCatalogItem
	err := r.db.QueryRow(ctx, query, packageID).Scan(&item.ID, &item.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) AggregatePackageMetrics(ctx context.Context, startDate, endDate time.Time) (map[string]packageMetric, error) {
	query := `
		SELECT
			o.package_id::text,
			COALESCE(p.name, 'Unknown Package') AS package_name,
			COUNT(*)::int AS total_transactions,
			COALESCE(SUM(o.total_price), 0)::numeric(18,2) AS total_revenue
		FROM subscription.orders o
		LEFT JOIN subscription.packages p ON p.id = o.package_id
		WHERE o.status = 'PAID'
		  AND o.package_id IS NOT NULL
		  AND o.created_at >= $1
		  AND o.created_at <= $2
		GROUP BY o.package_id, p.name
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]packageMetric)
	for rows.Next() {
		var item packageMetric
		if err := rows.Scan(&item.PackageID, &item.PackageName, &item.TotalTransactions, &item.TotalRevenue); err != nil {
			return nil, err
		}
		result[item.PackageID] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) AggregatePackageMonthlyMetrics(ctx context.Context, startDate, endDate time.Time) ([]packageMonthlyMetric, error) {
	query := `
		SELECT
			date_trunc('month', o.created_at) AS month_start,
			o.package_id::text,
			COALESCE(p.name, 'Unknown Package') AS package_name,
			COUNT(*)::int AS total_sales
		FROM subscription.orders o
		LEFT JOIN subscription.packages p ON p.id = o.package_id
		WHERE o.status = 'PAID'
		  AND o.package_id IS NOT NULL
		  AND o.created_at >= $1
		  AND o.created_at <= $2
		GROUP BY 1, 2, 3
		ORDER BY 1 ASC
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]packageMonthlyMetric, 0)
	for rows.Next() {
		var item packageMonthlyMetric
		if err := rows.Scan(&item.MonthStart, &item.PackageID, &item.PackageName, &item.TotalSales); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if items == nil {
		items = []packageMonthlyMetric{}
	}
	return items, nil
}

func (r *repository) AggregatePackageMonthlyMetricsByPackage(ctx context.Context, packageID string, startDate, endDate time.Time) (map[string]int, error) {
	query := `
		SELECT
			date_trunc('month', o.created_at) AS month_start,
			COUNT(*)::int AS total_sales
		FROM subscription.orders o
		WHERE o.status = 'PAID'
		  AND o.package_id IS NOT NULL
		  AND LOWER(o.package_id::text) = LOWER($1)
		  AND o.created_at >= $2
		  AND o.created_at <= $3
		GROUP BY 1
		ORDER BY 1 ASC
	`

	rows, err := r.db.Query(ctx, query, packageID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var monthStart time.Time
		var total int
		if err := rows.Scan(&monthStart, &total); err != nil {
			return nil, err
		}
		result[monthKey(monthStart)] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) GetPackageTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(o.total_price), 0)::numeric(18,2)
		FROM subscription.orders o
		WHERE o.status = 'PAID'
		  AND o.package_id IS NOT NULL
		  AND o.created_at >= $1
		  AND o.created_at <= $2
	`

	var total float64
	if err := r.db.QueryRow(ctx, query, startDate, endDate).Scan(&total); err != nil {
		return 0, fmt.Errorf("gagal mengambil total revenue: %w", err)
	}
	return total, nil
}

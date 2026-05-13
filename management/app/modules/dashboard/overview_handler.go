package dashboard

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func fetchInternal(url string, dest *map[string]any) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(dest)
}

// GetOverviewDashboard handles GET /api/superadmin/dashboard/overview
func (h *Handler) GetOverviewDashboard(c *gin.Context) {
	usersBase := os.Getenv("USERS_SERVICE_URL")
	if usersBase == "" {
		usersBase = "http://localhost:9002"
	}
	ticketsBase := os.Getenv("TICKETS_SERVICE_URL")
	if ticketsBase == "" {
		ticketsBase = "http://localhost:2004"
	}
	subscriptionBase := os.Getenv("SUBSCRIPTION_SERVICE_URL")
	if subscriptionBase == "" {
		subscriptionBase = "http://localhost:5004"
	}

	var (
		kycSummary  = map[string]any{}
		suppSummary = map[string]any{}
		opsSummary  = map[string]any{}
		kycErr      error
		suppErr     error
		opsErr      error
		wg          sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		kycErr = fetchInternal(usersBase+"/internal/dashboard/kyc-summary", &kycSummary)
	}()

	go func() {
		defer wg.Done()
		suppErr = fetchInternal(ticketsBase+"/internal/dashboard/support-summary", &suppSummary)
	}()

	go func() {
		defer wg.Done()
		opsErr = fetchInternal(subscriptionBase+"/internal/dashboard/ops-summary", &opsSummary)
	}()

	wg.Wait()

	if kycErr != nil {
		kycSummary = map[string]any{"error": kycErr.Error()}
	}
	if suppErr != nil {
		suppSummary = map[string]any{"error": suppErr.Error()}
	}
	if opsErr != nil {
		opsSummary = map[string]any{"error": opsErr.Error()}
	}

	c.JSON(http.StatusOK, gin.H{
		"kyc":          kycSummary,
		"support":      suppSummary,
		"operational":  opsSummary,
	})
}

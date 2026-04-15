package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
)

// ==========================================
// CONTROLLER STRUCT
// ==========================================

type Controller struct {
	dispatcher *concurrency.Dispatcher
	response   *ehttp.ResponseHelper
	service    Service
}

func NewController(dispatcher *concurrency.Dispatcher, response *ehttp.ResponseHelper, service Service) *Controller {
	return &Controller{
		dispatcher: dispatcher,
		response:   response,
		service:    service,
	}
}

// ==========================================
// HTTP HANDLERS
// ==========================================

// CreateOrderHandler — POST /api/orders (PBI-45)
// Hanya bisa diakses role CLIENT/USER (enforced di gateway routes.json)
func (c *Controller) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari header yang di-inject oleh gateway setelah validasi JWT
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		// Fallback: coba ambil dari Authorization Bearer JWT claims
		// (gateway seharusnya sudah inject X-User-ID, ini safety net)
		w.WriteHeader(http.StatusUnauthorized)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Anda harus login terlebih dahulu",
		})
		return
	}

	var dto CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Format request tidak valid",
		})
		return
	}

	// Dispatch ke concurrency worker (pola ZaFramework)
	payload := map[string]interface{}{
		"user_id": userID,
		"dto":     dto,
	}
	result, err := c.dispatcher.DispatchAndWait(r.Context(), "create_order", payload, concurrency.PriorityHigh)
	if err != nil {
		// Map error message ke status code yang sesuai PBI
		statusCode := http.StatusBadRequest
		msg := err.Error()

		if strings.Contains(msg, "tidak ditemukan") || strings.Contains(msg, "tidak tersedia") {
			statusCode = http.StatusBadRequest
		} else if strings.Contains(msg, "tidak teridentifikasi") {
			statusCode = http.StatusUnauthorized
		} else {
			statusCode = http.StatusInternalServerError
		}

		w.WriteHeader(statusCode)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": msg,
		})
		return
	}

	// 201 Created (PBI-45 AC: response 201 Created jika pesanan berhasil dibuat)
	w.WriteHeader(http.StatusCreated)
	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Pesanan berhasil dibuat",
		"data":    result,
	})
}

// ==========================================
// WORKER JOB PROCESSORS
// ==========================================

func (s *orderService) ProcessCreateOrderJobController(ctx context.Context, payload interface{}) (interface{}, error) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid payload type for CreateOrderJob")
	}
	userID, _ := data["user_id"].(string)
	dto, ok := data["dto"].(CreateOrderDTO)
	if !ok {
		return nil, fmt.Errorf("invalid dto type for CreateOrderJob")
	}
	return s.CreateOrder(ctx, userID, dto)
}

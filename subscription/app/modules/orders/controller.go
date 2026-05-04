package orders

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
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

func roleCode(r *http.Request) string {
	return strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
}

func extractOrderID(r *http.Request, prefix string) string {
	id := r.PathValue("id")
	if id != "" {
		return id
	}
	if strings.HasPrefix(r.URL.Path, prefix) {
		return strings.TrimPrefix(r.URL.Path, prefix)
	}
	return ""
}

// CreateOrderHandler — POST /api/orders (PBI-45)
func (c *Controller) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "CLIENT" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role CLIENT",
		})
		return
	}

	// Ambil user_id dari header yang di-inject oleh gateway setelah validasi JWT
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
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
	if dto.DurationMonths <= 0 {
		dto.DurationMonths = 1
	}

	// Dispatch ke concurrency worker (pola ZaFramework)
	payload := map[string]interface{}{
		"user_id": userID,
		"dto":     dto,
	}
	result, err := c.dispatcher.DispatchAndWait(r.Context(), "create_order", payload, concurrency.PriorityHigh)
	if err != nil {
		statusCode := http.StatusBadRequest
		msg := err.Error()
		msgLower := strings.ToLower(msg)

		if strings.Contains(msg, "tidak teridentifikasi") {
			statusCode = http.StatusUnauthorized
		} else if strings.HasPrefix(msgLower, "gagal") {
			statusCode = http.StatusInternalServerError
		} else if strings.Contains(msgLower, "kyc") {
			statusCode = http.StatusForbidden
		}

		w.WriteHeader(statusCode)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": msg,
		})
		return
	}

	order, ok := result.(*Order)
	if !ok {
		payloadBytes, _ := json.Marshal(result)
		var fallback Order
		if err := json.Unmarshal(payloadBytes, &fallback); err != nil || fallback.ID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": "gagal memproses respons pembuatan pesanan",
			})
			return
		}
		order = &fallback
	}

	// 201 Created (PBI-45 AC: response 201 Created jika pesanan berhasil dibuat)
	w.WriteHeader(http.StatusCreated)
	c.response.JSON(w, r, map[string]interface{}{
		"success":        true,
		"message":        "Pesanan berhasil dibuat",
		"order_id":       order.ID,
		"invoice_number": order.InvoiceNumber,
		"package_id":     order.PackageID,
		"total_price":    order.TotalPrice,
		"payment_method": order.PaymentMethod,
		"status":         order.Status,
		"created_at":     order.CreatedAt,
		"data":           order,
	})
}

// ListOrdersClientHandler — GET /api/orders
func (c *Controller) ListOrdersClientHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "CLIENT" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role CLIENT",
		})
		return
	}

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Anda harus login terlebih dahulu",
		})
		return
	}

	q := r.URL.Query()
	pageVal, _ := strconv.Atoi(q.Get("page"))
	limitVal, _ := strconv.Atoi(q.Get("limit"))
	filter := ClientOrderFilter{
		Status:    strings.TrimSpace(q.Get("status")),
		StartDate: strings.TrimSpace(q.Get("start_date")),
		EndDate:   strings.TrimSpace(q.Get("end_date")),
		Page:      pageVal,
		Limit:     limitVal,
	}

	list, err := c.service.ListOrdersForClient(r.Context(), userID, filter)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "gagal") {
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"data":    list,
	})
}

// GetOrderDetailClientHandler — GET /api/orders/{id}
func (c *Controller) GetOrderDetailClientHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "CLIENT" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role CLIENT",
		})
		return
	}

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Anda harus login terlebih dahulu",
		})
		return
	}

	orderID := extractOrderID(r, "/api/orders/")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	order, err := c.service.GetOrderDetailForClient(r.Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			w.WriteHeader(http.StatusNotFound)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": "Data pesanan tidak ditemukan",
			})
			return
		case errors.Is(err, ErrOrderForbidden):
			w.WriteHeader(http.StatusForbidden)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": "Anda tidak memiliki akses ke pesanan ini",
			})
			return
		default:
			statusCode := http.StatusBadRequest
			if strings.Contains(strings.ToLower(err.Error()), "gagal") {
				statusCode = http.StatusInternalServerError
			}
			w.WriteHeader(statusCode)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": err.Error(),
			})
			return
		}
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success":                   true,
		"order_id":                  order.OrderID,
		"invoice_number":            order.InvoiceNumber,
		"package_name":              order.PackageName,
		"total_price":               order.TotalPrice,
		"payment_method":            order.PaymentMethod,
		"status":                    order.Status,
		"has_payment_proof":         order.HasPaymentProof,
		"payment_proof_uploaded_at": order.PaymentProofUploadedAt,
		"payment_proof_url":         order.PaymentProofURL,
		"created_at":                order.CreatedAt,
		"data":                      order,
	})
}

// UploadPaymentProofClientHandler — POST /api/orders/{id}/payment-proof
func (c *Controller) UploadPaymentProofClientHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "CLIENT" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role CLIENT",
		})
		return
	}

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Anda harus login terlebih dahulu",
		})
		return
	}

	orderID := extractOrderID(r, "/api/orders/")
	orderID = strings.TrimSuffix(orderID, "/payment-proof")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	if err := r.ParseMultipartForm(6 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Format upload tidak valid",
		})
		return
	}

	file, header, err := r.FormFile("payment_proof")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "file bukti transfer wajib diunggah",
		})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, 5*1024*1024+1))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "gagal membaca file bukti transfer",
		})
		return
	}
	if len(data) > 5*1024*1024 {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ukuran file bukti transfer maksimal 5MB",
		})
		return
	}

	contentType := strings.TrimSpace(header.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	result, err := c.service.UploadPaymentProof(r.Context(), userID, orderID, PaymentProofFile{
		FileName:    header.Filename,
		ContentType: contentType,
		Data:        data,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, ErrOrderForbidden):
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success":                   true,
		"message":                   result.Message,
		"order_id":                  result.OrderID,
		"has_payment_proof":         result.HasPaymentProof,
		"payment_proof_uploaded_at": result.PaymentProofUploadedAt,
		"data":                      result,
	})
}

// GetPaymentProofClientHandler — GET /api/orders/{id}/payment-proof
func (c *Controller) GetPaymentProofClientHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "CLIENT" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role CLIENT",
		})
		return
	}

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Anda harus login terlebih dahulu",
		})
		return
	}

	orderID := extractOrderID(r, "/api/orders/")
	orderID = strings.TrimSuffix(orderID, "/payment-proof")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	proof, err := c.service.GetPaymentProofForClient(r.Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound), errors.Is(err, ErrPaymentProofMissing):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, ErrOrderForbidden):
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	contentType := strings.TrimSpace(proof.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "private, max-age=60")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(proof.Data)
}

// ListOrdersAdminHandler — GET /api/admin/orders
func (c *Controller) ListOrdersAdminHandler(w http.ResponseWriter, r *http.Request) {
	role := roleCode(r)
	if role != "OPERASIONAL" && role != "SUPERADMIN" && role != "CEO" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role OPERASIONAL",
		})
		return
	}

	q := r.URL.Query()
	pageVal, _ := strconv.Atoi(q.Get("page"))
	limitVal, _ := strconv.Atoi(q.Get("limit"))
	adminFilter := AdminOrderFilter{
		Status:    strings.TrimSpace(q.Get("status")),
		Search:    strings.TrimSpace(q.Get("search")),
		StartDate: strings.TrimSpace(q.Get("start_date")),
		EndDate:   strings.TrimSpace(q.Get("end_date")),
		Page:      pageVal,
		Limit:     limitVal,
	}

	list, err := c.service.ListOrdersForAdmin(r.Context(), adminFilter)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "gagal") {
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"data":    list,
	})
}

// GetOrderDetailAdminHandler — GET /api/admin/orders/{id}
func (c *Controller) GetOrderDetailAdminHandler(w http.ResponseWriter, r *http.Request) {
	role := roleCode(r)
	if role != "OPERASIONAL" && role != "SUPERADMIN" && role != "CEO" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role OPERASIONAL",
		})
		return
	}

	orderID := extractOrderID(r, "/api/admin/orders/")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	order, err := c.service.GetOrderDetailForAdmin(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": "Data pesanan tidak ditemukan",
			})
			return
		}

		statusCode := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "gagal") {
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"data":    order,
	})
}

// GetPaymentProofAdminHandler — GET /api/admin/orders/{id}/payment-proof
func (c *Controller) GetPaymentProofAdminHandler(w http.ResponseWriter, r *http.Request) {
	role := roleCode(r)
	if role != "OPERASIONAL" && role != "SUPERADMIN" && role != "CEO" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role OPERASIONAL",
		})
		return
	}

	orderID := extractOrderID(r, "/api/admin/orders/")
	orderID = strings.TrimSuffix(orderID, "/payment-proof")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	proof, err := c.service.GetPaymentProofForAdmin(r.Context(), orderID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound), errors.Is(err, ErrPaymentProofMissing):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	contentType := strings.TrimSpace(proof.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "private, max-age=60")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(proof.Data)
}

// VerifyOrderHandler — PATCH /api/admin/orders/{id}/verify
func (c *Controller) VerifyOrderHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "OPERASIONAL" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role OPERASIONAL",
		})
		return
	}

	orderID := extractOrderID(r, "/api/admin/orders/")
	orderID = strings.TrimSuffix(orderID, "/verify")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	var dto VerifyOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Format request tidak valid",
		})
		return
	}
	if strings.TrimSpace(dto.Action) == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "field action wajib diisi",
		})
		return
	}

	action := strings.ToUpper(strings.TrimSpace(dto.Action))
	if action == "REJECT" && strings.TrimSpace(dto.RejectReason) == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "alasan reject wajib diisi",
		})
		return
	}

	adminID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	result, err := c.service.VerifyOrder(r.Context(), orderID, action, dto.RejectReason, adminID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": "Data pesanan tidak ditemukan",
			})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success":           true,
		"message":           result.Message,
		"order_id":          result.OrderID,
		"new_status":        result.NewStatus,
		"verification_note": result.VerificationNote,
		"data":              result,
	})
}

// ActivateOrderHandler — PATCH /api/admin/orders/{id}/activate
// Endpoint ini dipakai untuk proses otomatis sistem.
func (c *Controller) ActivateOrderHandler(w http.ResponseWriter, r *http.Request) {
	role := roleCode(r)
	if role != "" && role != "SYSTEM" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "aktivasi hanya dapat diproses otomatis oleh sistem",
		})
		return
	}

	if strings.TrimSpace(r.Header.Get("X-System-Action")) != "AUTO_ACTIVATE" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "aktivasi hanya dapat diproses otomatis oleh sistem",
		})
		return
	}

	orderID := extractOrderID(r, "/api/admin/orders/")
	orderID = strings.TrimSuffix(orderID, "/activate")
	if orderID == "" {
		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "ID pesanan harus disertakan pada URL",
		})
		return
	}

	activation, err := c.service.ActivateOrderSystem(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			w.WriteHeader(http.StatusNotFound)
			c.response.JSON(w, r, map[string]interface{}{
				"success":       false,
				"error_message": "Data pesanan tidak ditemukan",
			})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success":         true,
		"message":         activation.Message,
		"subscription_id": activation.SubscriptionID,
		"start_date":      activation.StartDate,
		"end_date":        activation.EndDate,
		"status":          activation.Status,
		"data":            activation,
	})
}

// GetMySubscriptionHandler — GET /api/subscriptions/me
func (c *Controller) GetMySubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if roleCode(r) != "CLIENT" {
		w.WriteHeader(http.StatusForbidden)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "endpoint ini hanya dapat diakses oleh role CLIENT",
		})
		return
	}

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": "Anda harus login terlebih dahulu",
		})
		return
	}

	subs, err := c.service.GetActiveSubscriptions(r.Context(), userID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "gagal") {
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		c.response.JSON(w, r, map[string]interface{}{
			"success":       false,
			"error_message": err.Error(),
		})
		return
	}

	if len(subs) == 0 {
		c.response.JSON(w, r, map[string]interface{}{
			"success":                    true,
			"message":                    "Belum berlangganan",
			"data":                       nil,
			"active_subscriptions":       []SubscriptionStatus{},
			"total_active_subscriptions": 0,
		})
		return
	}

	// data tetap dipertahankan sebagai fallback kompatibilitas lama (single object).
	current := subs[0]
	c.response.JSON(w, r, map[string]interface{}{
		"success":                    true,
		"data":                       current,
		"active_subscriptions":       subs,
		"total_active_subscriptions": len(subs),
	})
}

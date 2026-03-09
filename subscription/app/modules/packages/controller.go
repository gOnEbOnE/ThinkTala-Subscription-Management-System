package packages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
)

// ==========================================
// 1. CONTROLLER STRUCT
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
// 2. HTTP HANDLERS DI BAWAH INI
// ==========================================

// CreatePackageHandler - POST /api/packages
func (c *Controller) CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePackageDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON body",
		})
		return
	}

	result, err := c.dispatcher.DispatchAndWait(r.Context(), "create_package", payload, concurrency.PriorityHigh)
	if err != nil {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Paket berhasil dibuat",
		"data":    result,
	})
}

// GetPackagesAdminHandler - GET /api/packages (With Query Params, Epic04 PBI-33, PBI-36)
func (c *Controller) GetPackagesAdminHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")

	payload := map[string]string{
		"status":    status,
		"min_price": minPrice,
		"max_price": maxPrice,
	}

	result, err := c.dispatcher.DispatchAndWait(r.Context(), "get_admin_packages", payload, concurrency.PriorityNormal)
	if err != nil {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Always return 200 OK even if array is empty
	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Berhasil mengambil data paket",
		"data":    result,
	})
}

// GetCatalogHandler - GET /api/packages/catalog (Epic04 PBI-37)
func (c *Controller) GetCatalogHandler(w http.ResponseWriter, r *http.Request) {
	result, err := c.dispatcher.DispatchAndWait(r.Context(), "get_catalog_packages", nil, concurrency.PriorityNormal)
	if err != nil {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Berhasil mengambil katalog belanja",
		"data":    result,
	})
}

// UpdatePackageHandler - PUT /account/packages/{id} (Epic04 PBI-34)
func (c *Controller) UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract UUID dari path: /api/admin/packages/{id}
	id := r.PathValue("id")
	if id == "" {
		// fallback: manual extract
		const prefix = "/api/admin/packages/"
		if len(r.URL.Path) > len(prefix) {
			id = r.URL.Path[len(prefix):]
		}
	}
	if id == "" {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": "ID paket harus disertakan pada URL",
		})
		return
	}

	var payload UpdatePackageDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": "Invalid JSON body",
		})
		return
	}

	// Bungkus id & payload ke map untuk request worker
	reqPayload := map[string]interface{}{
		"id":   id,
		"data": payload,
	}

	result, err := c.dispatcher.DispatchAndWait(r.Context(), "update_package", reqPayload, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "paket tidak ditemukan atau sudah dihapus" {
			c.response.JSON(w, r, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Paket berhasil diperbarui",
		"data":    result,
	})
}

// DeletePackageHandler - DELETE /account/packages/{id} (Epic04 PBI-35)
func (c *Controller) DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract UUID dari path: /api/admin/packages/{id}
	id := r.PathValue("id")
	if id == "" {
		// fallback: manual extract
		const prefix = "/api/admin/packages/"
		if len(r.URL.Path) > len(prefix) {
			id = r.URL.Path[len(prefix):]
		}
	}
	if id == "" {
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": "ID paket harus disertakan pada URL",
		})
		return
	}

	_, err := c.dispatcher.DispatchAndWait(r.Context(), "delete_package", id, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "paket tidak ditemukan atau sudah dihapus" {
			w.WriteHeader(http.StatusNotFound)
			c.response.JSON(w, r, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err.Error() == "tidak dapat menghapus paket yang sedang aktif" {
			w.WriteHeader(http.StatusBadRequest)
			c.response.JSON(w, r, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		c.response.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// PBI-35: Response successful delete message
	c.response.JSON(w, r, map[string]interface{}{
		"success": true,
		"message": "Paket berhasil dihapus",
	})
}

// ==========================================
// 3. WORKER JOB PROCESSORS (ZaFramework standard)
// ==========================================

func (s *packageService) ProcessCreatePackageJob(ctx context.Context, payload interface{}) (interface{}, error) {
	data, ok := payload.(CreatePackageDTO)
	if !ok {
		return nil, fmt.Errorf("invalid payload type for CreatePackageJob")
	}
	return s.CreatePackage(ctx, data)
}

func (s *packageService) ProcessGetAdminPackagesJob(ctx context.Context, payload interface{}) (interface{}, error) {
	data, ok := payload.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("invalid payload type for GetAdminPackagesJob")
	}
	return s.GetAdminPackages(ctx, data["status"], data["min_price"], data["max_price"])
}

func (s *packageService) ProcessGetCatalogJob(ctx context.Context, _ interface{}) (interface{}, error) {
	return s.GetCatalogPackages(ctx)
}

func (s *packageService) ProcessUpdatePackageJob(ctx context.Context, payload interface{}) (interface{}, error) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid payload type for UpdatePackageJob")
	}

	id := data["id"].(string)
	updatePayload := data["data"].(UpdatePackageDTO)

	return s.UpdatePackage(ctx, id, updatePayload)
}

func (s *packageService) ProcessDeletePackageJob(ctx context.Context, payload interface{}) (interface{}, error) {
	id, ok := payload.(string)
	if !ok {
		return nil, fmt.Errorf("invalid payload type for DeletePackageJob")
	}
	return nil, s.DeletePackage(ctx, id)
}

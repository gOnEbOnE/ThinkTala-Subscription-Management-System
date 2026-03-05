package users

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/master-abror/zaframework/core/concurrency"
	ehttp "github.com/master-abror/zaframework/core/http"
)

type Controller struct {
	Dispatcher *concurrency.Dispatcher
	Response   *ehttp.ResponseHelper
}

func NewController(d *concurrency.Dispatcher, r *ehttp.ResponseHelper) *Controller {
	return &Controller{
		Dispatcher: d,
		Response:   r,
	}
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	c.Response.View(w, r, "public/views/account/users/page.html", "Users Management", map[string]any{})
}

func (c *Controller) Inactive(w http.ResponseWriter, r *http.Request) {
	c.Response.View(w, r, "public/views/account/users/inactive.html", "Users Management", map[string]any{})
}

func (c *Controller) Active(w http.ResponseWriter, r *http.Request) {
	c.Response.View(w, r, "public/views/account/users/active.html", "Users Management", map[string]any{})
}

func (c *Controller) Banned(w http.ResponseWriter, r *http.Request) {
	c.Response.View(w, r, "public/views/account/users/banned.html", "Users Management", map[string]any{})
}

func (c *Controller) GetDetailUser(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID dari URL Path (Go 1.22+)
	id := r.PathValue("id")
	if id == "" {
		c.Response.JSON(w, r, map[string]any{
			"status": false,
			"msg":    "ID User diperlukan",
		})
		return
	}

	// 2. Siapkan Payload untuk Worker
	payload := map[string]any{
		"id": id,
	}

	// 3. Dispatch Job
	// Pastikan key "get_detail_user" sudah didaftarkan di Registry Worker
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "get_detail_user", payload, concurrency.PriorityHigh)

	if err != nil {
		c.Response.JSON(w, r, map[string]any{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}

	// 4. Casting Result (Result dari service adalah *User)
	user, ok := result.(*User)
	if !ok {
		c.Response.JSON(w, r, map[string]any{
			"status": false,
			"msg":    "Internal Error: Gagal casting data user",
		})
		return
	}

	// 5. Return Response
	c.Response.JSON(w, r, map[string]any{
		"status": true,
		"msg":    "Detail user berhasil diambil",
		"data":   user,
	})
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID dari URL (Go 1.22+)
	id := r.PathValue("id")
	if id == "" {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": "ID User diperlukan"})
		return
	}

	// 2. Parse Form Data (Pengganti JSON Decoder)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
        c.Response.JSON(w, r, map[string]any{"status": false, "msg": "Gagal memproses data multipart: " + err.Error()})
        return
    }

	// 3. Susun Payload Manual dari Form Value
	// r.FormValue mengambil data berdasarkan atribut 'name' di HTML
	input := map[string]any{
		"id":        id,
		"fullname":  r.FormValue("fullname"),
		"nik":       r.FormValue("nik"),
		"nrk":       r.FormValue("nrk"),
		"is_active": r.FormValue("is_active"), // Akan bernilai "1" atau "0"
	}

	fmt.Println("input", input)

	// 4. Dispatch Job
	_, err := c.Dispatcher.DispatchAndWait(r.Context(), "update_user", input, concurrency.PriorityHigh)

	if err != nil {
		c.Response.JSON(w, r, map[string]any{"status": false, "msg": err.Error()})
		return
	}

	c.Response.JSON(w, r, map[string]any{
		"status": true,
		"msg":    "Data berhasil diperbarui",
	})
}

func (c *Controller) GetInactiveUsers(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil Parameter dari URL
	// Default value ditangani nanti di Service, disini cukup parsing
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("q")

	// Konversi ke tipe yang sesuai (int)
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// 2. Siapkan Payload untuk Worker
	payload := map[string]any{
		"page":   page,
		"limit":  limit,
		"search": search,
	}

	// 3. Dispatch Job
	// Nama job: "get_inactive_users" (Harus didaftarkan di router worker nanti)
	// Priority: Normal (karena ini load data tabel, bukan login krusial)
	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "get_inactive_users", payload, concurrency.PriorityNormal)

	// 4. Error Handling dari Worker
	if err != nil {
		c.Response.JSON(w, r, map[string]any{
			"status": false,
			"msg":    err.Error(),
		})
		return
	}

	// 5. Casting Result
	// Service kita mengembalikan map[string]any (berisi "data" dan "meta")
	dataResult, ok := result.(map[string]any)
	if !ok {
		c.Response.JSON(w, r, map[string]any{
			"status": false,
			"msg":    "Internal Error: Invalid result format from worker",
		})
		return
	}

	// 6. Return Response
	// Kita tambahkan flag status: true agar konsisten dengan standar API Anda
	// Tapi key "data" dan "meta" tetap ada di root agar terbaca oleh Script Frontend
	response := map[string]any{
		"status": true,
		"msg":    "Data berhasil diambil",
		"data":   dataResult["data"],
		"meta":   dataResult["meta"],
	}

	c.Response.JSON(w, r, response)
}

package dashboard

import (
	"net/http"

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
	c.Response.View(w, r, "public/views/account/dashboard/page.html", "Welcome to ZAFramework", map[string]any{})
}

package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pintoter/warehouse-api/internal/model"
)

type getProductsResponse struct {
	Products []model.Product `json:"products"`
}

type errorResponse struct {
	Err string `json:"error"`
}

func renderJSON(w http.ResponseWriter, r *http.Request, code int, data any) {
	log.Printf("[Response] [%s] %s - Status code: [%d]", r.Method, r.URL.Path, code)

	resp, _ := json.MarshalIndent(data, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(resp)
}

package transport

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/pintoter/warehouse-api/internal/model"
// 	"github.com/pintoter/warehouse-api/pkg/logger"
// )

// func (h *Handler) reserveProducts(w http.ResponseWriter, r *http.Request) {
// 	var input reserveProductsReq
// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil || len(input.Products) == 0 {
// 		renderJSON(w, r, http.StatusBadRequest, errorResponse{model.ErrInvalidInput.Error()})
// 		return
// 	}

// 	products := h.service.ReserveProducts(r.Context(), input.Products)

// 	renderJSON(w, r, http.StatusCreated, products)
// }

// func (h *Handler) releaseProducts(w http.ResponseWriter, r *http.Request) {
// 	var input releaseProductsReq
// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil || len(input.Products) == 0 {
// 		renderJSON(w, r, http.StatusBadRequest, errorResponse{model.ErrInvalidInput.Error()})
// 		return
// 	}

// 	products := h.service.ReleaseProducts(r.Context(), input.Products)
// 	logger.DebugKV(r.Context(), "res", "products", products)

// 	renderJSON(w, r, http.StatusCreated, products)
// }

// func (h *Handler) showProducts(w http.ResponseWriter, r *http.Request) {
// 	var input showProductsReq
// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil || input.WarehouseId <= 0 {
// 		renderJSON(w, r, http.StatusBadRequest, errorResponse{model.ErrInvalidInput.Error()})
// 		return
// 	}

// 	products, err := h.service.GetProductsByWarehouse(r.Context(), input.WarehouseId)
// 	if err != nil {
// 		renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
// 		return
// 	}

// 	renderJSON(w, r, http.StatusCreated, getProductsResponse{Products: products});

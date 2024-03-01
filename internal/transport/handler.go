package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pintoter/warehouse-api/internal/service"
	"github.com/pintoter/warehouse-api/pkg/logger"
)

type Handler struct {
	router  *mux.Router
	service service.ProductService
}

func NewHandler(service service.ProductService) *Handler {
	handler := &Handler{
		router:  mux.NewRouter(),
		service: service,
	}

	handler.InitRoutes()

	return handler
}

func (h *Handler) InitRoutes() {
	h.router.HandleFunc("/reserveProducts", h.reserveProducts).Methods(http.MethodPost)
	h.router.HandleFunc("/releaseProducts", h.releaseProducts).Methods(http.MethodPost)
	h.router.HandleFunc("/showProducts", h.showProducts).Methods(http.MethodPost)
	// v1 := h.router.PathPrefix("/api/v1").Subrouter()
	// {
	// 	v1.HandleFunc("/reserveProducts", h.reserveProducts).Methods(http.MethodPost)
	// 	v1.HandleFunc("/realeseProducts", h.realeseProducts).Methods(http.MethodPost)
	// 	v1.HandleFunc("/showProducts", h.showProducts).Methods(http.MethodPost)
	// }
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.InfoKV(r.Context(), "New request", "Addr", r.RemoteAddr)
	h.router.ServeHTTP(w, r)
}

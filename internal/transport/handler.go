package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
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

	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterService(service, "ProductService")
	handler.router.Handle("/rpc", rpcServer)

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.InfoKV(r.Context(), "New request", "Addr", r.RemoteAddr)
	h.router.ServeHTTP(w, r)
}

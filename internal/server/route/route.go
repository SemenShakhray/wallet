package route

import (
	"wallet/internal/config"
	"wallet/internal/server/service"

	"github.com/go-chi/chi"
)

func NewRouter(cfg config.Config, serv service.Servicer) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Get("/api/v1/wallets/{wallet_id}", serv.GetBalance)
	r.Post("/api/v1/wallet", serv.Deposited)

	return r, nil
}

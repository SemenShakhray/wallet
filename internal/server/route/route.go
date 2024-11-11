package route

import (
	"wallet/internal/server/service"

	"github.com/go-chi/chi"
)

func NewRouter(serv service.Servicer) (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Get("/api/v1/wallets/{wallet_id}", serv.GetBalance)
	r.Post("/api/v1/wallet", serv.Deposited)

	return r, nil
}

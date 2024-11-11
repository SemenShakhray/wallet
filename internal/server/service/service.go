package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"wallet/internal/config"
	"wallet/internal/models"
	"wallet/internal/server/repository"

	"github.com/google/uuid"
)

type Servicer interface {
	GetBalance(w http.ResponseWriter, r *http.Request)
	Deposited(w http.ResponseWriter, r *http.Request)
	CloseConnectionDB()
}

type Service struct {
	Repository repository.Repositorer
}

func NewService(cfg config.Config) (Servicer, error) {
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		return nil, err
	}
	return &Service{
		Repository: repo,
	}, nil
}

func (s *Service) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parts := strings.Split(r.URL.Path, "/")
	wallet_id := parts[len(parts)-1]
	_, err := uuid.Parse(wallet_id)
	if err != nil {
		http.Error(w, `{"error":"wrong wallet uuid"}`, http.StatusBadRequest)
		return
	}
	ok, err := s.Repository.CheckExist(ctx, wallet_id)
	if err != nil {
		http.Error(w, `{"error":"check exist is failed"}`, http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, `{"error":"wallet does not exist"}`, http.StatusNotFound)
		return
	}

	balance, err := s.Repository.GetBalance(ctx, wallet_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("wallet_id: %s\nbalance: %d", wallet_id, balance)))
}

func (s *Service) Deposited(w http.ResponseWriter, r *http.Request) {
	wallet := models.Wallet{}

	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		http.Error(w, `{"error":"incorrect deserialization JSON"}`, http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(wallet.UUID)
	if err != nil {
		http.Error(w, `{"error":"wrong wallet uuid"}`, http.StatusBadRequest)
		return
	}

	ok, err := s.Repository.CheckExist(ctx, wallet.UUID)
	if err != nil {
		http.Error(w, `{"error":"check exist is failed"}`, http.StatusInternalServerError)
		return
	}

	if !ok {
		err := s.Repository.CreteWallet(ctx, wallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("wallet_id: %s\nbalance: %d", wallet.UUID, wallet.Amount)))
		return
	}

	balance, err := s.Repository.Deposited(ctx, wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("wallet_id: %s\nbalance: %d", wallet.UUID, balance)))
}

func (s *Service) CloseConnectionDB() {
	s.Repository.CloseConnectionDB()
}

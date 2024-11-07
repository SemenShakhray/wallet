package repository

import (
	"context"
	"wallet/internal/config"
	"wallet/internal/models"
	"wallet/internal/storage"
)

type Repositorer interface {
	CheckExist(ctx context.Context, uuid string) (bool, error)
	GetBalance(ctx context.Context, uuid string) (int, error)
	CreteWallet(ctx context.Context, wallet models.Wallet) error
	Deposited(ctx context.Context, wallet models.Wallet) (int, error)
	CloseConnectionDB()
}

type Repository struct {
	Store storage.Storer
}

func NewRepository(cfg config.Config) (Repositorer, error) {
	store, err := storage.ConnectDB(cfg.DSN)
	if err != nil {
		return nil, err
	}

	return Repository{Store: store}, nil
}

func (repo Repository) CheckExist(ctx context.Context, uuid string) (bool, error) {
	return repo.Store.CheckExist(ctx, uuid)
}

func (repo Repository) GetBalance(ctx context.Context, uuid string) (int, error) {
	return repo.Store.GetBalance(ctx, uuid)
}

func (repo Repository) CreteWallet(ctx context.Context, wallet models.Wallet) error {
	return repo.Store.CreteWallet(ctx, wallet)
}

func (repo Repository) Deposited(ctx context.Context, wallet models.Wallet) (int, error) {
	return repo.Store.Deposited(ctx, wallet)
}

func (repo Repository) CloseConnectionDB() {
	repo.Store.CloseConnectionDB()
}

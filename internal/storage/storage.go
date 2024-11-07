package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"wallet/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storer interface {
	GetBalance(ctx context.Context, uuid string) (int, error)
	CheckExist(ctx context.Context, uuid string) (bool, error)
	CreteWallet(ctx context.Context, wallet models.Wallet) error
	Deposited(ctx context.Context, wallet models.Wallet) (int, error)
	CloseConnectionDB()
}

type Store struct {
	db *sql.DB
}

func ConnectDB(dsn string) (Storer, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed connected database: %v", err)
	}

	//create table for wallet

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wallet (
	id SERIAL PRIMARY KEY,
	wallet_id VARCHAR(128),
	amount INT)`)

	if err != nil {
		return nil, err
	}

	return Store{
		db: db,
	}, nil
}

func (store Store) CheckExist(ctx context.Context, uuid string) (bool, error) {
	var idEx int
	row := store.db.QueryRowContext(ctx, `SELECT id FROM wallet WHERE wallet_id=$1`, uuid)
	err := row.Scan(&idEx)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s Store) GetBalance(ctx context.Context, uuid string) (int, error) {
	w := models.Wallet{}
	row := s.db.QueryRowContext(ctx, `SELECT amount FROM wallet WHERE wallet_id=$1`, uuid)
	err := row.Scan(&w.Amount)
	if err != nil {
		return 0, err
	}
	return w.Amount, nil
}

func (s Store) CreteWallet(ctx context.Context, wallet models.Wallet) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO wallet (wallet_id, amount) VALUES ($1, $2)`, wallet.UUID, wallet.Amount)
	log.Println(wallet.UUID, wallet.Amount)
	if err != nil {
		return fmt.Errorf(`{"error":"add wallet"}`)
	}
	return nil
}

func (s Store) Deposited(ctx context.Context, wallet models.Wallet) (int, error) {
	w := models.Wallet{}
	err := s.db.QueryRow(`SELECT amount FROM wallet WHERE wallet_id=$1`, wallet.UUID).Scan(&w.Amount)
	if err != nil {
		return w.Amount, err
	}
	w.Amount = w.Amount + wallet.Amount
	_, err = s.db.ExecContext(ctx, `UPDATE wallet SET amount = $1 WHERE wallet_id = $2`, w.Amount, wallet.UUID)
	if err != nil {
		return w.Amount, fmt.Errorf(`{"error":"couldn't update the balance"}`)
	}
	return w.Amount, nil
}

func (s Store) CloseConnectionDB() {
	s.db.Close()
}

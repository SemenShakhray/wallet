package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	DB *sql.DB
}

func ConnectDB(dsn string) (Storer, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed connected database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wallet (
	id SERIAL PRIMARY KEY,
	wallet_id VARCHAR(128),
	amount INT);
	CREATE INDEX IF NOT EXISTS wallet_uuid ON wallet (wallet_id);`)

	if err != nil {
		return nil, err
	}

	return Store{
		DB: db,
	}, nil
}

func (store Store) CheckExist(ctx context.Context, uuid string) (bool, error) {
	var idEx int
	row := store.DB.QueryRowContext(ctx, `SELECT id FROM wallet WHERE wallet_id=$1`, uuid)
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
	row := s.DB.QueryRowContext(ctx, `SELECT amount FROM wallet WHERE wallet_id=$1`, uuid)
	err := row.Scan(&w.Amount)
	if err != nil {
		return 0, err
	}
	return w.Amount, nil
}

func (s Store) CreteWallet(ctx context.Context, wallet models.Wallet) error {
	switch strings.ToLower(wallet.OpperationType) {
	case "deposit":
	case "withdraw":
		wallet.Amount = -wallet.Amount
	}
	_, err := s.DB.ExecContext(ctx, `INSERT INTO wallet (wallet_id, amount) VALUES ($1, $2)`, wallet.UUID, wallet.Amount)
	log.Println(wallet.UUID, wallet.Amount)
	if err != nil {
		return fmt.Errorf(`{"error":"add wallet"}`)
	}
	return nil
}

func (s Store) Deposited(ctx context.Context, wallet models.Wallet) (int, error) {
	w := models.Wallet{}
	err := s.DB.QueryRow(`SELECT amount FROM wallet WHERE wallet_id = $1`, wallet.UUID).Scan(&w.Amount)
	if err != nil {
		return w.Amount, err
	}
	switch strings.ToLower(wallet.OpperationType) {
	case "deposit":
		w.Amount = w.Amount + wallet.Amount
	case "withdraw":
		w.Amount = w.Amount - wallet.Amount
	default:
		return w.Amount, fmt.Errorf(`{"error":"failed opperation with wallet"}`)
	}

	_, err = s.DB.ExecContext(ctx, `UPDATE wallet SET amount = $1 WHERE wallet_id = $2`, w.Amount, wallet.UUID)
	if err != nil {
		return w.Amount, fmt.Errorf(`{"error":"couldn't update the balance"}`)
	}
	return w.Amount, nil
}

func (s Store) CloseConnectionDB() {
	s.DB.Close()
}

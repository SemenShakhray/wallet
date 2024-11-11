package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"wallet/internal/models"
	"wallet/internal/store"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := store.Store{DB: db}
	ctx := context.Background()

	t.Run("wallet exists", func(t *testing.T) {
		walletID := "550e8400-e29b-41d4-a716-446655440000"

		mock.ExpectQuery("SELECT id FROM wallet WHERE wallet_id=?").
			WithArgs(walletID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		exists, err := store.CheckExist(ctx, walletID)
		assert.NoError(t, err)
		if !exists {
			assert.Error(t, fmt.Errorf("expected wallet to exist, but it does not"))
		}
	})

	t.Run("wallet does not exist", func(t *testing.T) {
		walletID := "550e8400-e29b-41d4-a716-446655440000"

		mock.ExpectQuery("SELECT id FROM wallet WHERE wallet_id=?").
			WithArgs(walletID).
			WillReturnError(sql.ErrNoRows)

		exists, err := store.CheckExist(ctx, walletID)
		assert.NoError(t, err)
		if !exists {
			assert.Error(t, fmt.Errorf("expected wallet to exist, but it does not"))
		}
	})
}

func TestGetBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := store.Store{DB: db}
	ctx := context.Background()

	t.Run("valid balance", func(t *testing.T) {
		walletID := "550e8400-e29b-41d4-a716-446655440000"
		expectedBalance := 100

		mock.ExpectQuery("SELECT amount FROM wallet WHERE wallet_id=?").
			WithArgs(walletID).
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(expectedBalance))

		balance, err := store.GetBalance(ctx, walletID)
		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, balance)
	})
}

func TestCreteWallet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := store.Store{DB: db}
	ctx := context.Background()

	t.Run("create wallet successfully", func(t *testing.T) {
		wallet := models.Wallet{
			UUID:           "550e8400-e29b-41d4-a716-446655440000",
			Amount:         100,
			OpperationType: "deposit",
		}

		mock.ExpectExec("INSERT INTO wallet").
			WithArgs(wallet.UUID, wallet.Amount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := store.CreteWallet(ctx, wallet)
		assert.NoError(t, err)
	})
}

func TestDeposited(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := store.Store{DB: db}
	ctx := context.Background()

	t.Run("deposit amount successfully", func(t *testing.T) {
		wallet := models.Wallet{
			UUID:           "550e8400-e29b-41d4-a716-446655440000",
			Amount:         100,
			OpperationType: "deposit",
		}

		mock.ExpectQuery(`SELECT amount FROM wallet WHERE wallet_id = \$1`).
			WithArgs(wallet.UUID).
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(200))

		mock.ExpectExec(`UPDATE wallet SET amount = \$1 WHERE wallet_id = \$2`).
			WithArgs(300, wallet.UUID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		newBalance, err := store.Deposited(ctx, wallet)
		assert.NoError(t, err)
		assert.Equal(t, 300, newBalance)
	})

	t.Run("withdraw amount successfully", func(t *testing.T) {
		wallet := models.Wallet{
			UUID:           "550e8400-e29b-41d4-a716-446655440000",
			Amount:         50,
			OpperationType: "withdraw",
		}

		mock.ExpectQuery(`SELECT amount FROM wallet WHERE wallet_id = ?`).
			WithArgs(wallet.UUID).
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(200))

		mock.ExpectExec(`UPDATE wallet SET amount = \$1 WHERE wallet_id = \$2`).
			WithArgs(150, wallet.UUID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		newBalance, err := store.Deposited(ctx, wallet)
		assert.NoError(t, err)
		assert.Equal(t, 150, newBalance)
	})
}

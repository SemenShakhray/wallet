package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet/internal/models"
	"wallet/internal/server/repository"
	"wallet/internal/server/service"
	mock_store "wallet/internal/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockStore := mock_store.NewMockStorer(ctrl)
	repo := repository.Repository{Store: mockStore}
	service := service.Service{Repository: repo}

	t.Run("OK GetBalance", func(t *testing.T) {
		walletID := "550e8400-e29b-41d4-a716-446655440000"
		balance := 200

		mockStore.EXPECT().CheckExist(ctx, walletID).Return(true, nil)
		mockStore.EXPECT().GetBalance(ctx, walletID).Return(200, nil)

		req := httptest.NewRequest("GET", "/wallets/"+walletID, nil)
		rec := httptest.NewRecorder()

		service.GetBalance(rec, req)
		res := rec.Body.String()
		assert.Equal(t, res, fmt.Sprintf("wallet_id: %s\nbalance: %d", walletID, balance))
	})

	t.Run("Wallet doesn't exist", func(t *testing.T) {

		walletID := uuid.NewString()
		mockStore.EXPECT().CheckExist(ctx, walletID).Return(false, nil)

		req := httptest.NewRequest(http.MethodGet, "/wallets/"+walletID, nil)
		rec := httptest.NewRecorder()

		service.GetBalance(rec, req)

		status := rec.Result().StatusCode
		assert.Equal(t, 404, status)
	})

	t.Run("Invalid UUID", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/wallets/"+"fsdfsdfsd", nil)
		rec := httptest.NewRecorder()

		service.GetBalance(rec, req)

		status := rec.Result().StatusCode
		assert.Equal(t, 400, status)
	})
}

func TestDeposit(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockStore := mock_store.NewMockStorer(ctrl)
	repo := repository.Repository{Store: mockStore}
	service := service.Service{Repository: repo}

	t.Run("OK_CreateWallet", func(t *testing.T) {

		uuid := "c6c9be7e-d7bf-4828-aeac-539d786034ed"
		wallet := models.Wallet{
			UUID:           uuid,
			OpperationType: "Deposit",
			Amount:         100,
		}

		body := fmt.Sprintf(`{
		"uuid":           "%s",
			"opperationType": "Deposit",
			"amount":         100
			}`, uuid)

		mockStore.EXPECT().CheckExist(ctx, uuid).Return(false, nil)
		mockStore.EXPECT().CreteWallet(ctx, wallet).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader(body))
		rec := httptest.NewRecorder()

		service.Deposited(rec, req)

		status := rec.Result().StatusCode
		assert.Equal(t, 200, status)

		res := rec.Body.String()
		assert.Equal(t, res, fmt.Sprintf("wallet_id: %s\nbalance: %d", uuid, wallet.Amount))
	})

	t.Run("OK Deposit", func(t *testing.T) {

		uuid := "c6c9be7e-d7bf-4828-aeac-539d786034ed"
		wallet := models.Wallet{
			UUID:           uuid,
			OpperationType: "Deposit",
			Amount:         100,
		}

		body := fmt.Sprintf(`{
		"uuid":           "%s",
			"opperationType": "Deposit",
			"amount":         100
			}`, uuid)

		mockStore.EXPECT().CheckExist(ctx, wallet.UUID).Return(true, nil)
		mockStore.EXPECT().Deposited(ctx, wallet).Return(100, nil)

		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader(body))
		rec := httptest.NewRecorder()

		service.Deposited(rec, req)

		status := rec.Result().StatusCode
		assert.Equal(t, 200, status)

		res := rec.Body.String()
		assert.Equal(t, res, fmt.Sprintf("wallet_id: %s\nbalance: %d", wallet.UUID, 100))
	})

}

package service_test

import (
	"scheduler/internal/app/models"
	"scheduler/internal/app/repository"
	"scheduler/internal/app/service"
	mockrepo "scheduler/tests/unittest/mock/repo"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func NewMockAccountService(mockRepo repository.AccountRepo) service.AccountAuthService {
	return &service.AccountService{
		Repo: mockRepo,
	}
}

func TestLogin(t *testing.T) {
	mockRepo := &mockrepo.MockAccountRepo{}
	mockRequest := &models.LoginRequest{
		Email:    "staff123@gmail.com",
		Password: "staff123",
	}

	mockAccount := &models.Account{
		ID:       1,
		Name:     "staff123",
		Email:    "staff123@gmail.com",
		Phone:    "123456789",
		Password: "$2a$12$Z3f60AnCWhJSnJsrdw7yae8IBjE/VXHYUHvC/NN0UOhLEO.BaPyJW",
		Role:     "staff",
	}

	mockRepo.LoginFunc = func(*models.LoginRequest) (*models.Account, error) {
		return mockAccount, nil
	}

	c, _ := gin.CreateTestContext(nil)

	accountService := NewMockAccountService(mockRepo)
	account, err := accountService.LoginAccount(c, mockRequest)

	t.Run("test normal login path", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, mockAccount, account)
	})
}

func TestCreateAccount(t *testing.T) {
	mockRepo := &mockrepo.MockAccountRepo{}
	mockRequest := &models.CreateAccountRequest{
		Email:    "staff123@gmail.com",
		Password: "Staff@1234",
		Role:     "staff",
		Phone:    "123456789",
		Name:     "Staff 123",
	}

	mockAccount := &models.Account{
		ID:    1,
		Email: "staff123@gmail.com",
		Role:  "staff",
		Phone: "123456789",
		Name:  "Staff 123",
	}

	mockRepo.CreateAccountFunc = func(*models.Account) (*models.Account, error) {
		return mockAccount, nil
	}

	c, _ := gin.CreateTestContext(nil)

	accountService := NewMockAccountService(mockRepo)
	account, err := accountService.CreateAccount(c, mockRequest)

	t.Run("test normal account creation", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, mockAccount, account)
	})
	

	// Check wrong password error
	mockRequest = &models.CreateAccountRequest{
		Email:    "staff123@gmail.com",
		Password: "Staff234",
		Role:     "staff",
		Phone:    "123456789",
		Name:     "Staff 123",
	}
	_, err = accountService.CreateAccount(c, mockRequest)
	t.Run("test password error", func(t *testing.T) {
		assert.Error(t, err)
	})
}

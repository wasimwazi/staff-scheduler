package mockservice

import (
	"scheduler/internal/app/models"

	"github.com/gin-gonic/gin"
)

type MockAccountService struct {
	LoginFunc         func(*models.LoginRequest) (*models.Account, error)
	CreateAccountFunc func(*models.CreateAccountRequest) (*models.AccountResponse, error)
}

func (m *MockAccountService) LoginAccount(c *gin.Context,
	req *models.LoginRequest) (*models.Account, error) {

	return m.LoginFunc(req)
}

func (m *MockAccountService) CreateAccount(c *gin.Context,
	req *models.CreateAccountRequest) (*models.AccountResponse, error) {

	return m.CreateAccountFunc(req)
}

package mockrepo

import (
	"scheduler/internal/app/models"

	"github.com/gin-gonic/gin"
)

type MockAccountRepo struct {
	LoginFunc         func(*models.LoginRequest) (*models.Account, error)
	CreateAccountFunc func(*models.Account) (*models.Account, error)
}

func (m *MockAccountRepo) Login(c *gin.Context, user *models.LoginRequest) (*models.Account, error) {
	return m.LoginFunc(user)
}

func (m *MockAccountRepo) CreateAccount(c *gin.Context, account *models.Account) (*models.Account, error) {
	return m.CreateAccountFunc(account)
}

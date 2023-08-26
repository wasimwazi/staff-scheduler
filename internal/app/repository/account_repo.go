package repository

import (
	"errors"
	"scheduler/internal/app/models"
	"scheduler/internal/platform/database"
	"scheduler/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewAccountRepo : Returns Account Repo
func NewAccountRepo(db *database.Database) AccountRepo {
	return &PostgresRepo{
		DB: db.Client,
	}
}

// Login : Postgres function to validate account crederntials
func (pg *PostgresRepo) Login(ctx *gin.Context, login *models.LoginRequest) (*models.Account, error) {
	var account models.Account
	if err := pg.DB.Where("email = ?", login.Email).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (pg *PostgresRepo) CreateAccount(ctx *gin.Context, account *models.Account) (*models.Account, error) {
	var existingAccount models.Account
	err := pg.DB.Where("email = ?", account.Email).First(&existingAccount).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err == nil {
		return nil, utils.ErrEmailAlreadyExists
	}

	err = pg.DB.Create(&account).Error
	if err != nil {
		return nil, err
	}

	return account, nil
}

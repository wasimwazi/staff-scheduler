package service

import (
	"errors"
	"scheduler/internal/app/models"
	repo "scheduler/internal/app/repository"
	"scheduler/internal/platform/database"
	"scheduler/utils"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AccountService : Account Service Struct
type AccountService struct {
	Repo repo.AccountRepo
}

// NewAccountService : Returns Account Service
func NewAccountService(db *database.Database) AccountAuthService {
	return &AccountService{
		Repo: repo.NewAccountRepo(db),
	}
}

// LoginAccount : Login to the account
func (us *AccountService) LoginAccount(ctx *gin.Context, loginRequest *models.LoginRequest) (*models.Account, error) {
	loggedInAccount, err := us.Repo.Login(ctx, loginRequest)
	if err != nil {
		return nil, err
	}

	storedPassword := []byte(loggedInAccount.Password)
	if err := bcrypt.CompareHashAndPassword(storedPassword, []byte(loginRequest.Password)); err != nil {
		return nil, errors.New(utils.ErrPassword)
	}

	return loggedInAccount, nil
}

func (s *AccountService) CreateAccount(ctx *gin.Context, createRequest *models.CreateAccountRequest) (*models.AccountResponse, error) {
	if err := CheckPasswordStrength(createRequest.Password); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createRequest.Password = string(hashedPassword)

	newAccount := models.Account{
		Name:     createRequest.Name,
		Email:    createRequest.Email,
		Phone:    createRequest.Phone,
		Password: createRequest.Password,
		Role:     createRequest.Role,
	}

	createdAccount, err := s.Repo.CreateAccount(ctx, &newAccount)
	if err != nil {
		return nil, err
	}

	respAccount := &models.AccountResponse{
		ID:    createdAccount.ID,
		Name:  createRequest.Name,
		Email: createdAccount.Email,
		Phone: createRequest.Phone,
		Role:  createRequest.Role,
	}

	return respAccount, nil
}

func CheckPasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUppercase   bool
		hasLowercase   bool
		hasDigit       bool
		hasSpecialChar bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	if !(hasUppercase && hasLowercase && hasDigit && hasSpecialChar) {
		return utils.ErrPasswordCondition
	}

	return nil
}

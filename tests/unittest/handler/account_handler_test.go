package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scheduler/internal/app/handler"
	"scheduler/internal/app/models"
	"scheduler/internal/app/service"
	mockservice "scheduler/tests/unittest/mock/service"
	mockutils "scheduler/tests/unittest/mock/utils"
	"scheduler/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	router := gin.Default()

	mockInput := `{
		"email": "staff123@gmail.com",
		"password": "staff123"
	}`

	mockAccountService := new(mockservice.MockAccountService)
	mockAccountService.LoginFunc = func(loginReq *models.LoginRequest) (*models.Account, error) {
		return &models.Account{}, nil
	}

	mockAuth := new(mockutils.MockTokenGenerator)
	mockAuth.GenerateTokenFunc = func(id uint, email string, role string) (*models.JWT, error) {
		jwt := &models.JWT{
			ID:          id,
			AccessToken: "mockAccesToken",
		}
		return jwt, nil
	}

	handler := &handler.Handler{
		Services: &service.Services{
			AccountService: mockAccountService,
		},
		JWTGenerator: mockAuth,
	}
	router.POST("/login", handler.LoginAccount)

	reqBody1 := strings.NewReader(mockInput)
	req, _ := http.NewRequest(http.MethodPost, "/login", reqBody1)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w1.Body.Bytes(), &responseBody)

	message, ok := responseBody["message"].(string)
	if !ok {
		t.Fatal("error decoding the message field from response body")
	}
	t.Run("test login normal path", func(t *testing.T) {
		assert.True(t, ok)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w1.Code)
		assert.Equal(t, utils.MsgLoggedIn, message)
	})

	// Test account not found

	mockAccountService.LoginFunc = func(loginReq *models.LoginRequest) (*models.Account, error) {
		return nil, gorm.ErrRecordNotFound
	}

	reqBody2 := strings.NewReader(mockInput)
	req, _ = http.NewRequest(http.MethodPost, "/login", reqBody2)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	t.Run("test login invalid user path", func(t *testing.T) {
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, w2.Code)
	})

	// Test inactive user

	mockAccountService.LoginFunc = func(loginReq *models.LoginRequest) (*models.Account, error) {
		return nil, utils.ErrUserInactive
	}

	reqBody3 := strings.NewReader(mockInput)
	req, _ = http.NewRequest(http.MethodPost, "/login", reqBody3)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req)
	t.Run("test login inactive user", func(t *testing.T) {
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, w3.Code)
	})

	// Test input validation

	mockAccountService.LoginFunc = func(loginReq *models.LoginRequest) (*models.Account, error) {
		return nil, utils.ErrUserInactive
	}

	mockInput = `{
		"email": "",
		"password": "staff123"
	}
	`
	reqBody4 := strings.NewReader(mockInput)
	req, _ = http.NewRequest(http.MethodPost, "/login", reqBody4)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req)
	t.Run("test login invalid input", func(t *testing.T) {
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, w4.Code)
	})
}

func TestCreateAccount(t *testing.T) {
	router := gin.Default()

	mockInput1 := `{
		"email": "staff123@gmail.com",
		"password": "Staff#123",
		"role": "staff",
		"name": "staff 123",
		"phone": "+91123456789"
	}`

	accountResp := &models.AccountResponse{
		ID:    utils.UserID1,
		Email: utils.StaffEmail123,
		Role:  utils.RoleStaff,
		Name:  utils.StaffName123,
		Phone: utils.StaffPhone,
	}

	mockAccountService := new(mockservice.MockAccountService)
	mockAccountService.CreateAccountFunc = func(loginReq *models.CreateAccountRequest) (*models.AccountResponse, error) {
		return accountResp, nil
	}

	handler := &handler.Handler{
		Services: &service.Services{
			AccountService: mockAccountService,
		},
	}
	router.POST("/createaccount", handler.CreateAccount)

	reqBody1 := strings.NewReader(mockInput1)
	req, _ := http.NewRequest(http.MethodPost, "/createaccount", reqBody1)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)

	var response struct {
		Message string                 `json:"message"`
		Payload models.AccountResponse `json:"payload"`
	}
	err := json.Unmarshal(w1.Body.Bytes(), &response)
	assert.NoError(t, err)

	t.Run("test account creation normal path", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w1.Code)
		assert.Equal(t, utils.MsgAccountCreated, response.Message)
		assert.Equal(t, utils.UserID1, response.Payload.ID)
		assert.Equal(t, utils.StaffName123, response.Payload.Name)
		assert.Equal(t, utils.StaffEmail123, response.Payload.Email)
		assert.Equal(t, utils.StaffPhone, response.Payload.Phone)
		assert.Equal(t, utils.RoleStaff, response.Payload.Role)
	})

	// Test invalid role field

	mockInput2 := `{
		"email": "staff123@gmail.com",
		"password": "Staff#123",
		"role": "teacher",
		"name": "staff 123",
		"phone": "+91123456789"
	}`

	reqBody2 := strings.NewReader(mockInput2)
	req, _ = http.NewRequest(http.MethodPost, "/createaccount", reqBody2)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	var errResponse utils.ErrResponse
	err = json.Unmarshal(w2.Body.Bytes(), &errResponse)
	assert.NoError(t, err)

	t.Run("test account creation invalid role", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, w2.Code)
		assert.Equal(t, utils.ErrInvalidRole.Error(), errResponse.Error)
	})
}

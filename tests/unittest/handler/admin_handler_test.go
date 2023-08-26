package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scheduler/internal/app/handler"
	"scheduler/internal/app/models"
	"scheduler/internal/app/service"
	mockservice "scheduler/tests/unittest/mock/service"
	"scheduler/utils"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestEditUser(t *testing.T) {
	router := gin.Default()

	mockInput := `{
		"name": "staff-edited123",
		"phone": "+91123456789"
	}`

	accountResp := models.AccountResponse{
		ID:    utils.UserID1,
		Name:  utils.StaffNameEdited,
		Phone: utils.StaffPhone,
		Email: utils.StaffEmail123,
		Role:  utils.RoleStaff,
	}

	mockAdminService := new(mockservice.MockAdminService)
	mockAdminService.EditUserFunc = func(userID uint, req *models.EditUserRequest) (*models.AccountResponse, error) {
		return &accountResp, nil
	}
	handler := &handler.Handler{
		Services: &service.Services{
			AdminService: mockAdminService,
		},
	}
	router.PUT("/users/:user_id", handler.EditUserHandler)

	reqBody1 := strings.NewReader(mockInput)
	req, _ := http.NewRequest(http.MethodPut, "/users/1", reqBody1)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)

	var response struct {
		Message string                 `json:"message"`
		Payload models.AccountResponse `json:"payload"`
	}
	err := json.Unmarshal(w1.Body.Bytes(), &response)
	assert.NoError(t, err)

	t.Run("test edit user normal path", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w1.Code)
		assert.Equal(t, utils.MsgUserEditedSuccess, response.Message)
		assert.Equal(t, utils.StaffEmail123, response.Payload.Email)
		assert.Equal(t, utils.StaffNameEdited, response.Payload.Name)
		assert.Equal(t, utils.StaffPhone, response.Payload.Phone)
		assert.Equal(t, utils.UserID1, response.Payload.ID)
		assert.Equal(t, utils.RoleStaff, response.Payload.Role)
	})

	mockAdminService.EditUserFunc = func(userID uint, req *models.EditUserRequest) (*models.AccountResponse, error) {
		return &models.AccountResponse{}, gorm.ErrRecordNotFound
	}

	reqBody2 := strings.NewReader(mockInput)
	req2, _ := http.NewRequest(http.MethodPut, "/users/1", reqBody2)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var errResponse utils.ErrResponse
	err = json.Unmarshal(w2.Body.Bytes(), &errResponse)
	assert.NoError(t, err)

	t.Run("test edit user service error", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, w2.Code)
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), errResponse.Error)
	})
}

func ContextMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("currentUserID", userID)
		c.Next()
	}
}

func TestDeleteUser(t *testing.T) {
	router := gin.Default()
	router.Use(ContextMiddleware(uint(123)))

	mockAdminService := new(mockservice.MockAdminService)
	mockAdminService.DeleteUserFunc = func(userID uint) error {
		return nil
	}

	handler := &handler.Handler{
		Services: &service.Services{
			AdminService: mockAdminService,
		},
	}
	router.DELETE("/users/:user_id", handler.DeleteUserHandler)

	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)

	t.Run("test delete user normal path", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w1.Code)
	})

	// Test the service returning error case
	mockAdminService.DeleteUserFunc = func(userID uint) error {
		return gorm.ErrRecordNotFound
	}

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)

	var errResp utils.ErrResponse
	err := json.Unmarshal(w2.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	t.Run("test delete user service error", func(t *testing.T) {
		assert.Equal(t, http.StatusNotFound, w2.Code)
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), errResp.Error)
	})
}

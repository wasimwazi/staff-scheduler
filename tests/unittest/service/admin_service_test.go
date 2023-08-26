package service_test

import (
	"scheduler/internal/app/models"
	"scheduler/internal/app/repository"
	"scheduler/internal/app/service"
	mockrepo "scheduler/tests/unittest/mock/repo"
	"scheduler/utils"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func NewMockAdminService(mockRepo repository.AdminRepo) service.AdminService {
	return &service.Admin{
		Repo: mockRepo,
	}
}

func TestEditUser(t *testing.T) {
	mockRepo := &mockrepo.MockAdminRepo{}

	mockRequest := &models.EditUserRequest{
		Name: utils.StaffNameEdited,
		Role: "admin",
	}
	userID := uint(1)

	mockAccount := &models.Account{
		ID:    1,
		Name:  utils.StaffNameEdited,
		Email: "staff123@gmail.com",
		Phone: "123456789",
		Role:  "admin",
	}

	expectedResponse := &models.AccountResponse{
		ID:    1,
		Name:  utils.StaffNameEdited,
		Email: "staff123@gmail.com",
		Phone: "123456789",
		Role:  "admin",
	}

	mockRepo.EditUserFunc = func(uint, *models.EditUserRequest) (*models.Account, error) {
		return mockAccount, nil
	}

	c, _ := gin.CreateTestContext(nil)

	adminService := NewMockAdminService(mockRepo)
	account, err := adminService.EditUser(c, userID, mockRequest)

	t.Run("test normal edit user", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, account)
	})

	mockRepo.EditUserFunc = func(uint, *models.EditUserRequest) (*models.Account, error) {
		return nil, utils.ErrFetchData
	}
	account, err = adminService.EditUser(c, userID, mockRequest)
	t.Run("test edit user repo error case", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, err, utils.ErrFetchData)
		assert.Equal(t, (*models.AccountResponse)(nil), account)
	})
}

func TestCreateSchedule(t *testing.T) {
	mockRepo := &mockrepo.MockAdminRepo{}

	mockRequest := &models.CreateScheduleRequest{
		UserID:         1,
		ShiftLength:    8,
		WorkDateString: "2023-08-23",
	}
	workDate, _ := time.Parse(utils.DateFormat, mockRequest.WorkDateString)

	mockSchedule := models.Schedule{
		ID:          1,
		WorkDate:    workDate,
		ShiftLength: 8,
		AccountID:   1,
		User:        utils.StaffName123,
	}

	expectedResponse := &models.CreateScheduleResponse{
		ID:          1,
		WorkDate:    workDate,
		ShiftLength: 8,
		UserID:      1,
		Username:    utils.StaffName123,
	}

	mockRepo.CreateScheduleFunc = func(*models.Schedule) (*models.Schedule, error) {
		return &mockSchedule, nil
	}

	mockUserIDNameRoleMap := models.UserIDNameRoleMap{
		1: models.NameRole{
			Name: utils.StaffName123,
			Role: utils.RoleStaff,
		},
	}
	mockRepo.GetUserIDNameRoleMapFunc = func([]uint) (models.UserIDNameRoleMap, error) {
		return mockUserIDNameRoleMap, nil
	}

	c, _ := gin.CreateTestContext(nil)

	adminService := NewMockAdminService(mockRepo)
	scheduleResp, err := adminService.CreateSchedule(c, mockRequest)

	t.Run("test normal create schedule path", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, scheduleResp)
	})
}

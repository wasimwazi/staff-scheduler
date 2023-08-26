package mockservice

import (
	"scheduler/internal/app/models"

	"github.com/gin-gonic/gin"
)

type MockAdminService struct {
	EditUserFunc                   func(uint, *models.EditUserRequest) (*models.AccountResponse, error)
	DeleteUserFunc                 func(uint) error
	CreateScheduleFunc       func(*models.CreateScheduleRequest) (*models.CreateScheduleResponse, error)
	CreateSchedulesFunc            func(*models.ScheduleRequest) ([]uint, []uint, error)
	EditScheduleFunc               func(uint, uint, *models.EditScheduleRequest) (*models.ScheduleResponse, error)
	DeleteScheduleFunc             func(uint) error
	GetUsersByAccumulatedHoursFunc func(string, string) ([]models.AccountWithAccumulatedShiftResponse, error)
}

func (m *MockAdminService) EditUser(c *gin.Context,
	userID uint, req *models.EditUserRequest) (*models.AccountResponse, error) {

	return m.EditUserFunc(userID, req)
}

func (m *MockAdminService) DeleteUser(c *gin.Context, userID uint) error {
	return m.DeleteUserFunc(userID)
}

func (m *MockAdminService) CreateSchedule(c *gin.Context, req *models.CreateScheduleRequest) (*models.CreateScheduleResponse, error) {
	return m.CreateScheduleFunc(req)
}

func (m *MockAdminService) EditSchedule(c *gin.Context, scheduleID uint,
	userID uint, req *models.EditScheduleRequest) (*models.ScheduleResponse, error) {

	return m.EditScheduleFunc(scheduleID, userID, req)
}

func (m *MockAdminService) DeleteSchedule(c *gin.Context, scheduleID uint) error {
	return m.DeleteScheduleFunc(scheduleID)
}

func (m *MockAdminService) GetUsersByAccumulatedHours(c *gin.Context, startDate,
	endDate string) ([]models.AccountWithAccumulatedShiftResponse, error) {

	return m.GetUsersByAccumulatedHoursFunc(startDate, endDate)
}

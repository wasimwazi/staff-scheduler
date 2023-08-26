package mockrepo

import (
	"scheduler/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
)

type MockAdminRepo struct {
	EditUserFunc                   func(uint, *models.EditUserRequest) (*models.Account, error)
	DeleteUserFunc                 func(uint) error
	GetUserIDNameRoleMapFunc       func([]uint) (models.UserIDNameRoleMap, error)
	CreateScheduleFunc       func(*models.Schedule) (*models.Schedule, error)
	CreateSchedulesFunc            func([]models.Schedule) error
	EditScheduleFunc               func(uint, uint, *models.EditScheduleRequest) (*models.Schedule, error)
	DeleteScheduleFunc             func(uint) error
	GetUsersByAccumulatedHoursFunc func(time.Time, time.Time) ([]models.AccountWithAccumulatedShiftResponse, error)
}

func (m *MockAdminRepo) EditUser(c *gin.Context, userID uint, req *models.EditUserRequest) (*models.Account, error) {
	return m.EditUserFunc(userID, req)
}

func (m *MockAdminRepo) DeleteUser(c *gin.Context, userID uint) error {
	return m.DeleteUserFunc(userID)
}

func (m *MockAdminRepo) GetUserIDNameRoleMap(c *gin.Context, userIDs []uint) (models.UserIDNameRoleMap, error) {
	return m.GetUserIDNameRoleMapFunc(userIDs)
}

func (m *MockAdminRepo) CreateSchedule(c *gin.Context, schedule *models.Schedule) (*models.Schedule, error) {
	return m.CreateScheduleFunc(schedule)
}

func (m *MockAdminRepo) CreateSchedules(c *gin.Context, schedules []models.Schedule) error {
	return m.CreateSchedulesFunc(schedules)
}

func (m *MockAdminRepo) EditSchedule(c *gin.Context, scheduleID uint, userID uint, editReq *models.EditScheduleRequest) (*models.Schedule, error) {
	return m.EditScheduleFunc(scheduleID, userID, editReq)
}

func (m *MockAdminRepo) DeleteSchedule(c *gin.Context, scheduleID uint) error {
	return m.DeleteScheduleFunc(scheduleID)
}

func (m *MockAdminRepo) GetUsersByAccumulatedHours(c *gin.Context, startDate time.Time, endDate time.Time) ([]models.AccountWithAccumulatedShiftResponse, error) {
	return m.GetUsersByAccumulatedHoursFunc(startDate, endDate)
}

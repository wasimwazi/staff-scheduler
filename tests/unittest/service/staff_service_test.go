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

func NewMockStaffService(mockRepo repository.StaffRepo) service.StaffService {
	return &service.Staff{
		Repo: mockRepo,
	}
}

func TestGetStaffSchedules(t *testing.T) {
	mockRepo := &mockrepo.MockStaffRepo{}

	mockSchedules := []models.UserScheduleResponse{
		{
			WorkDate:    "2023-08-24",
			ShiftLength: 8,
		},
		{
			WorkDate:    "2023-08-25",
			ShiftLength: 6,
		},
	}

	expectedResp := &models.GetScheduleResponse{
		UserID:    utils.UserID1,
		Schedules: mockSchedules,
	}

	mockRepo.GetStaffSchedulesFunc = func(uint, time.Time, time.Time) ([]models.UserScheduleResponse, error) {
		return mockSchedules, nil
	}

	c, _ := gin.CreateTestContext(nil)

	staffService := NewMockStaffService(mockRepo)
	userSchedules, err := staffService.GetStaffSchedules(c, utils.UserID1, utils.StartDate, utils.EndDate)

	t.Run("test normal GetStaffSchedules path", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, expectedResp, userSchedules)
	})

	mockRepo.GetStaffSchedulesFunc = func(uint, time.Time, time.Time) ([]models.UserScheduleResponse, error) {
		return nil, utils.ErrFetchData
	}
	userSchedules, err = staffService.GetStaffSchedules(c, utils.UserID1, utils.StartDate, utils.EndDate)
	t.Run("test error in GetStaffSchedules DB call", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, err, utils.ErrFetchData)
		assert.Equal(t, (*models.GetScheduleResponse)(nil), userSchedules)
	})

	mockRepo.GetStaffSchedulesFunc = func(uint, time.Time, time.Time) ([]models.UserScheduleResponse, error) {
		return mockSchedules, nil
	}
	userSchedules, err = staffService.GetStaffSchedules(c, utils.UserID1, "22-22-22", utils.EndDate)
	t.Run("test invalid date", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, err, utils.ErrInvalidDateInput)
		assert.Equal(t, (*models.GetScheduleResponse)(nil), userSchedules)
	})
}

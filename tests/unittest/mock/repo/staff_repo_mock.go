package mockrepo

import (
	"scheduler/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
)

type MockStaffRepo struct {
	GetStaffSchedulesFunc func(uint, time.Time, time.Time) ([]models.UserScheduleResponse, error)
}

func (m *MockStaffRepo) GetStaffSchedules(c *gin.Context, userID uint, startDate time.Time,
	endDate time.Time) ([]models.UserScheduleResponse, error) {
	return m.GetStaffSchedulesFunc(userID, startDate, endDate)
}

package mockservice

import (
	"scheduler/internal/app/models"

	"github.com/gin-gonic/gin"
)

type MockStaffService struct {
	GetStaffSchedulesFunc func(uint, string, string) (*models.GetScheduleResponse, error)
}

func (m *MockStaffService) GetStaffSchedules(c *gin.Context, userID uint, startDate string,
	endDate string) (*models.GetScheduleResponse, error) {

	return m.GetStaffSchedulesFunc(userID, startDate, endDate)
}

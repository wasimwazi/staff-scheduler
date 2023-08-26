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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetStaffSchedules(t *testing.T) {
	router := gin.Default()
	router.Use(ContextMiddleware(uint(123)))

	scheduleResp1 := models.GetScheduleResponse{
		UserID: uint(123),
		Schedules: []models.UserScheduleResponse{
			{
				WorkDate:    "2023-08-21",
				ShiftLength: 5,
			},
			{
				WorkDate:    "2023-08-22",
				ShiftLength: 7,
			},
		},
	}

	mockStaffService := new(mockservice.MockStaffService)
	mockStaffService.GetStaffSchedulesFunc = func(userID uint, startDate, endDate string) (*models.GetScheduleResponse, error) {
		return &scheduleResp1, nil
	}

	handler := &handler.Handler{
		Services: &service.Services{
			StaffService: mockStaffService,
		},
	}
	router.GET("/schedule", handler.GetScheduleHandler)

	req, _ := http.NewRequest(http.MethodGet, "/schedule?start_date=2023-08-21&end_date=2023-08-23", nil)

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)

	var response1 struct {
		Message string                     `json:"message"`
		Payload models.GetScheduleResponse `json:"payload"`
	}

	err := json.Unmarshal(w1.Body.Bytes(), &response1)
	assert.NoError(t, err)

	t.Run("test get schedule normal path", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w1.Code)
		assert.Equal(t, uint(123), response1.Payload.UserID)
		assert.Equal(t, "2023-08-21", response1.Payload.Schedules[0].WorkDate)
		assert.Equal(t, 5, response1.Payload.Schedules[0].ShiftLength)
		assert.Equal(t, "2023-08-22", response1.Payload.Schedules[1].WorkDate)
		assert.Equal(t, 7, response1.Payload.Schedules[1].ShiftLength)
	})

	// Test coworker schedule

	scheduleResp2 := models.GetScheduleResponse{
		UserID: utils.UserID1,
		Schedules: []models.UserScheduleResponse{
			{
				WorkDate:    "2023-08-22",
				ShiftLength: 6,
			},
		},
	}

	mockStaffService.GetStaffSchedulesFunc = func(userID uint, startDate, endDate string) (*models.GetScheduleResponse, error) {
		return &scheduleResp2, nil
	}

	req2, _ := http.NewRequest(http.MethodGet, "/schedule?staff_id=1&start_date=2023-08-21&end_date=2023-08-23", nil)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var response2 struct {
		Message string                     `json:"message"`
		Payload models.GetScheduleResponse `json:"payload"`
	}

	err = json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err)

	t.Run("test get coworker schedule normal path", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w2.Code)
		assert.Equal(t, utils.UserID1, response2.Payload.UserID)
		assert.Equal(t, "2023-08-22", response2.Payload.Schedules[0].WorkDate)
		assert.Equal(t, 6, response2.Payload.Schedules[0].ShiftLength)
	})
}

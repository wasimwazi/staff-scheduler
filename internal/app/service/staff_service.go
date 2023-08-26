package service

import (
	"scheduler/internal/app/models"
	repo "scheduler/internal/app/repository"
	"scheduler/internal/platform/database"
	"scheduler/utils"

	"github.com/gin-gonic/gin"
)

type Staff struct {
	Repo repo.StaffRepo
}

func NewStaffService(db *database.Database) StaffService {
	return &Staff{
		Repo: repo.NewStaffRepo(db),
	}
}

func (s *Staff) GetStaffSchedules(ctx *gin.Context, userID uint,
	startDateString, endDateString string) (*models.GetScheduleResponse, error) {

	startDate, endDate, err := utils.DateParser(startDateString, endDateString)
	if err != nil {
		return nil, err
	}

	schedules, err := s.Repo.GetStaffSchedules(ctx, userID, *startDate, *endDate)
	if err != nil {
		return nil, err
	}
	scheduleResponse := models.GetScheduleResponse{
		UserID:    userID,
		Schedules: schedules,
	}
	return &scheduleResponse, nil
}

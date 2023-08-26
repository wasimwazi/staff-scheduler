package repository

import (
	"scheduler/internal/app/models"
	"scheduler/internal/platform/database"
	"time"

	"github.com/gin-gonic/gin"
)

func NewStaffRepo(db *database.Database) StaffRepo {
	return &PostgresRepo{
		DB: db.Client,
	}
}

func (pg *PostgresRepo) GetStaffSchedules(ctx *gin.Context, userID uint, startDate time.Time,
	endDate time.Time) ([]models.UserScheduleResponse, error) {

	var userSchedules []models.UserScheduleResponse

	result := pg.DB.Table("schedules").
		Select("work_date, shift_length").
		Where("account_id = ? AND work_date BETWEEN ? AND ?", userID, startDate, endDate).
		Scan(&userSchedules)

	if result.Error != nil {
		return nil, result.Error
	}

	return userSchedules, nil
}

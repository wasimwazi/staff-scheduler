package repository

import (
	"errors"
	"scheduler/internal/app/models"
	"scheduler/internal/platform/database"
	"scheduler/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAdminRepo(db *database.Database) AdminRepo {
	return &PostgresRepo{
		DB: db.Client,
	}
}

func (pg *PostgresRepo) EditUser(c *gin.Context, userID uint, editRequest *models.EditUserRequest) (*models.Account, error) {
	var account models.Account
	if err := pg.DB.WithContext(c).First(&account, userID).Error; err != nil {
		return nil, err
	}

	// Update user fields if they are provided in the request
	if editRequest.Name != "" {
		account.Name = editRequest.Name
	}
	if editRequest.Phone != "" {
		account.Phone = editRequest.Phone
	}
	if editRequest.Role != "" {
		account.Role = editRequest.Role
	}

	if err := pg.DB.WithContext(c).Save(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (pg *PostgresRepo) DeleteUser(c *gin.Context, userID uint) error {
	var account models.Account
	if err := pg.DB.WithContext(c).First(&account, userID).Error; err != nil {
		return err
	}

	if account.DeletedAt.Valid {
		return utils.ErrUserInactive
	}
	if err := pg.DB.WithContext(c).Delete(&account).Error; err != nil {
		return err
	}

	return nil
}

func (pg *PostgresRepo) CreateSchedule(c *gin.Context, schedule *models.Schedule) (*models.Schedule, error) {
	var existingSchedule models.Schedule
	err := pg.DB.Where("account_id = ? AND work_date = ?", schedule.AccountID, schedule.WorkDate).First(&existingSchedule).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// If no existing schedule is found, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := pg.DB.Create(schedule).Error; err != nil {
			if database.IsUniqueConstraintViolation(err) {
				err = utils.ErrScheduleAlreadyExist
			}
			return nil, err
		}
		return schedule, nil
	}

	return nil, utils.ErrScheduleAlreadyExist
}

func (pg *PostgresRepo) GetUserIDNameRoleMap(c *gin.Context, userIds []uint) (models.UserIDNameRoleMap, error) {
	var users []models.User

	if err := pg.DB.Find(&users, userIds).Error; err != nil {
		return nil, err
	}

	userIDNameRoleMap := make(models.UserIDNameRoleMap)
	for _, user := range users {
		userIDNameRoleMap[user.ID] = models.NameRole{
			Name: user.Name,
			Role: user.Role,
		}
	}

	return userIDNameRoleMap, nil
}

func (pg *PostgresRepo) EditSchedule(c *gin.Context, scheduleID, userID uint,
	editRequest *models.EditScheduleRequest) (*models.Schedule, error) {

	var schedule models.Schedule
	if err := pg.DB.Where("id = ? AND account_id = ?", scheduleID, userID).First(&schedule).Error; err != nil {
		return nil, err
	}

	if editRequest.WorkDateString != "" {
		schedule.WorkDate = editRequest.WorkDate
	}
	if editRequest.ShiftLength != 0 {
		schedule.ShiftLength = editRequest.ShiftLength
	}

	if err := pg.DB.Save(&schedule).Error; err != nil {
		if database.IsUniqueConstraintViolation(err) {
			err = utils.ErrScheduleAlreadyExist
		}
		return nil, err
	}

	return &schedule, nil
}

func (pg *PostgresRepo) DeleteSchedule(c *gin.Context, scheduleID uint) error {
	var schedule models.Schedule
	if err := pg.DB.First(&schedule, scheduleID).Error; err != nil {
		return err
	}

	if err := pg.DB.Delete(&schedule).Error; err != nil {
		return err
	}

	return nil
}

func (pg *PostgresRepo) GetUsersByAccumulatedHours(c *gin.Context, startDate,
	endDate time.Time) ([]models.AccountWithAccumulatedShiftResponse, error) {

	var users []models.AccountWithAccumulatedShiftResponse

	result := pg.DB.Table("accounts").
		Select("accounts.id, accounts.email, accounts.name, accounts.phone, COALESCE(SUM(schedules.shift_length), 0) AS accumulated_work_hours").
		Joins("LEFT JOIN schedules ON accounts.id = schedules.account_id").
		Where("accounts.role = ? AND schedules.work_date BETWEEN ? AND ?", "staff", startDate, endDate).
		Group("accounts.id").
		Order("accumulated_work_hours DESC").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

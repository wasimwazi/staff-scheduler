package repository

import (
	"scheduler/internal/app/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	DB *gorm.DB
}

//Repo : Account repo interface
type AccountRepo interface {
	Login(*gin.Context, *models.LoginRequest) (*models.Account, error)
	CreateAccount(*gin.Context, *models.Account) (*models.Account, error)
}

type AdminRepo interface {
	EditUser(*gin.Context, uint, *models.EditUserRequest) (*models.Account, error)
	DeleteUser(*gin.Context, uint) error
	GetUserIDNameRoleMap(*gin.Context, []uint) (models.UserIDNameRoleMap, error)
	CreateSchedule(*gin.Context, *models.Schedule) (*models.Schedule, error)
	EditSchedule(*gin.Context, uint, uint, *models.EditScheduleRequest) (*models.Schedule, error)
	DeleteSchedule(*gin.Context, uint) error
	GetUsersByAccumulatedHours(*gin.Context, time.Time, time.Time) ([]models.AccountWithAccumulatedShiftResponse, error)
}

type StaffRepo interface {
	GetStaffSchedules(*gin.Context, uint, time.Time, time.Time) ([]models.UserScheduleResponse, error)
}

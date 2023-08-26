package service

import (
	"scheduler/internal/app/models"

	"github.com/gin-gonic/gin"
)

// AccountAuthService : Service interface for account
type AccountAuthService interface {
	LoginAccount(*gin.Context, *models.LoginRequest) (*models.Account, error)
	CreateAccount(ctx *gin.Context, account *models.CreateAccountRequest) (*models.AccountResponse, error)
}

// AdminService : Service interface for admin
type AdminService interface {
	EditUser(*gin.Context, uint, *models.EditUserRequest) (*models.AccountResponse, error)
	DeleteUser(*gin.Context, uint) error
	CreateSchedule(*gin.Context, *models.CreateScheduleRequest) (*models.CreateScheduleResponse, error)
	EditSchedule(*gin.Context, uint, uint, *models.EditScheduleRequest) (*models.ScheduleResponse, error)
	DeleteSchedule(*gin.Context, uint) error
	GetUsersByAccumulatedHours(*gin.Context, string, string) ([]models.AccountWithAccumulatedShiftResponse, error)
}

// StaffService : Service interface for staff
type StaffService interface {
	GetStaffSchedules(*gin.Context, uint, string, string) (*models.GetScheduleResponse, error)
}

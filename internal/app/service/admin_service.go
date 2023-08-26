package service

import (
	"fmt"
	"scheduler/internal/app/models"
	repo "scheduler/internal/app/repository"
	"scheduler/internal/platform/database"
	"scheduler/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type Admin struct {
	Repo repo.AdminRepo
}

func NewAdminService(db *database.Database) AdminService {
	return &Admin{
		Repo: repo.NewAdminRepo(db),
	}
}

func (ad *Admin) EditUser(c *gin.Context, userID uint, editRequest *models.EditUserRequest) (*models.AccountResponse, error) {
	account, err := ad.Repo.EditUser(c, userID, editRequest)
	if err != nil {
		return nil, err
	}
	accountResp := models.AccountResponse{
		Name:  account.Name,
		Email: account.Email,
		Phone: account.Phone,
		Role:  account.Role,
		ID:    account.ID,
	}
	return &accountResp, nil
}

func (ad *Admin) DeleteUser(c *gin.Context, userID uint) error {
	err := ad.Repo.DeleteUser(c, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ad *Admin) CreateSchedule(ctx *gin.Context, createRequest *models.CreateScheduleRequest) (*models.CreateScheduleResponse, error) {
	userID := createRequest.UserID
	userIDNameRoleMap, err := ad.Repo.GetUserIDNameRoleMap(ctx, []uint{userID})
	if err != nil {
		return nil, err
	}

	workDate, err := time.Parse(utils.DateFormat, createRequest.WorkDateString)
	if err != nil {
		return nil, utils.ErrInvalidDateInput
	}
	newSchedule := models.Schedule{
		User:        userIDNameRoleMap[userID].Name,
		WorkDate:    workDate,
		AccountID:   userID,
		ShiftLength: createRequest.ShiftLength,
	}

	fmt.Println("User role", userIDNameRoleMap[userID].Role)
	fmt.Println("User map", userIDNameRoleMap)

	if userIDNameRoleMap[userID].Role == utils.RoleStaff {
		schedule, err := ad.Repo.CreateSchedule(ctx, &newSchedule)
		if err != nil {
			return nil, err
		}

		scheduleResp := models.CreateScheduleResponse{
			ID:          schedule.ID,
			UserID:      schedule.AccountID,
			Username:    schedule.User,
			WorkDate:    schedule.WorkDate,
			ShiftLength: schedule.ShiftLength,
		}
		return &scheduleResp, nil
	}
	return nil, utils.ErrCreateScheduleForOtherUsers
}

func IsWorkDateUnique(schedules []models.ScheduleModel) bool {
	seenWorkDates := make(map[string]bool)

	for _, schedule := range schedules {
		if seenWorkDates[schedule.WorkDateString] {
			return false
		}
		seenWorkDates[schedule.WorkDateString] = true
	}

	return true
}

func (ad *Admin) EditSchedule(c *gin.Context, scheduleID, userID uint, editRequest *models.EditScheduleRequest) (*models.ScheduleResponse, error) {
	workDate, err := time.Parse(utils.DateFormat, editRequest.WorkDateString)
	if err != nil {
		return nil, err
	}
	editRequest.WorkDate = workDate
	schedule, err := ad.Repo.EditSchedule(c, scheduleID, userID, editRequest)
	if err != nil {
		return nil, err
	}

	formattedDate := schedule.WorkDate.Format("2006-01-02")
	if err != nil {
		return nil, err
	}
	scheduleResp := models.ScheduleResponse{
		ID:          schedule.ID,
		Name:        schedule.User,
		WorkDate:    formattedDate,
		ShiftLength: schedule.ShiftLength,
	}
	return &scheduleResp, nil
}

func (ad *Admin) DeleteSchedule(c *gin.Context, scheduleID uint) error {
	err := ad.Repo.DeleteSchedule(c, scheduleID)
	if err != nil {
		return err
	}

	return nil
}

func (ad *Admin) GetUsersByAccumulatedHours(c *gin.Context,
	startDateString, endDateString string) ([]models.AccountWithAccumulatedShiftResponse, error) {

	startDate, endDate, err := utils.DateParser(startDateString, endDateString)
	if err != nil {
		return nil, err
	}

	users, err := ad.Repo.GetUsersByAccumulatedHours(c, *startDate, *endDate)
	if err != nil {
		return nil, err
	}
	return users, nil
}

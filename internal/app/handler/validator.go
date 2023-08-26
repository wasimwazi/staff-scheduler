package handler

import (
	"scheduler/internal/app/models"
	"scheduler/utils"

	"github.com/go-playground/validator/v10"
)

func AccountRequestValidator(req interface{}, role string) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}

	if role != "" && role != utils.RoleAdmin && role != utils.RoleStaff {
		return utils.ErrInvalidRole
	}
	return nil
}

// Validate the provided json to check if the required fields are present or not
func validateParsedData(scheduleReq models.ScheduleRequest) bool {
	for _, user := range scheduleReq.Users {
		if user.UserID == 0 || user.Schedules == nil {
			return false
		}
		for _, schedule := range user.Schedules {
			if schedule.WorkDateString == "" || schedule.ShiftLength == 0 {
				return false
			}
		}
	}
	return true
}

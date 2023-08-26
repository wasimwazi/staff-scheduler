package handler

import (
	"errors"
	"net/http"
	"scheduler/internal/app/models"
	"scheduler/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (h *Handler) EditUserHandler(c *gin.Context) {
	var editRequest models.EditUserRequest
	if err := c.ShouldBindJSON(&editRequest); err != nil {
		log.Error("Error in EditUserHandler - ShouldBindJSON()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := AccountRequestValidator(editRequest, editRequest.Role); err != nil {
		log.Error("Error in EditUserHandler - AccountRequestValidator()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if editRequest.Name == "" && editRequest.Phone == "" && editRequest.Role == "" {
		log.Error("Error in EditUserHandler - atleat one field required")
		utils.Fail(c, http.StatusBadRequest, utils.ErrFieldMissing)
		return
	}

	userID, err := utils.ExtractUserIDFromContextParam(c)
	if err != nil {
		log.Error("Error in EditUserHandler - ExtractUserIDFromContextParam()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.Services.AdminService.EditUser(c, userID, &editRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error in EditUserHandler - record not found", err.Error())
			utils.Fail(c, http.StatusNotFound, gorm.ErrRecordNotFound.Error())
			return
		}
		log.Error("Error in EditUserHandler - unknown wrror from service", err.Error())
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp := Response{
		Message: utils.MsgUserEditedSuccess,
		Payload: user,
	}
	log.Info(utils.MsgUserEditedSuccess)
	utils.Send(c, http.StatusOK, resp)
}

func (h *Handler) DeleteUserHandler(c *gin.Context) {
	userID, err := utils.ExtractUserIDFromContextParam(c)
	if err != nil {
		log.Error("Error in DeleteUserHandler - ExtractUserIDFromContextParam()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	currentUser, err := utils.ExtractCurrentUserIDFromContext(c)
	if err != nil {
		log.Error("Error in DeleteUserHandler - ExtractCurrentUserIDFromContext()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if currentUser == userID {
		log.Error("Error in DeleteUserHandler - current user can't be deleted", utils.ErrCantDeletUser)
		utils.Fail(c, http.StatusForbidden, utils.ErrCantDeletUser)
		return
	}

	err = h.Services.AdminService.DeleteUser(c, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Error in DeleteUserHandler - user not found", err.Error())
		utils.Fail(c, http.StatusNotFound, gorm.ErrRecordNotFound.Error())
			return
		}
		log.Error("Error in DeleteUserHandler - unknown error from service", err.Error())
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp := Response{
		Message: utils.MsgUserDeleted,
		Payload: gin.H{
			"user_id": userID,
		},
	}
	log.Info(utils.MsgUserDeleted)
	utils.Send(c, http.StatusOK, resp)
}

func (h *Handler) CreateScheduleHandler(c *gin.Context) {
	var request models.CreateScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Error in CreateScheduleHandler - ShouldBindJSON()", err.Error())
		utils.Fail(c, http.StatusBadRequest, utils.ErrBadInput)
		return
	}

	userID, err := utils.ExtractUserIDFromContextParam(c)
	if err != nil {
		log.Error("Error in CreateScheduleHandler - ExtractUserIDFromContextParam()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	request.UserID = userID

	schedule, err := h.Services.AdminService.CreateSchedule(c, &request)
	if err != nil {
		log.Error("Error in CreateScheduleHandler - service CreateSchedule", err.Error())
		utils.Fail(c, http.StatusInternalServerError, "Failed to create schedule for the user")
		return
	}

	msg := utils.MsgScheduleCreated
	resp := Response{
		Message: msg,
		Payload: schedule,
	}
	log.Info(utils.MsgScheduleCreated)
	utils.Send(c, http.StatusCreated, resp)
}

func (h *Handler) EditScheduleHandler(c *gin.Context) {
	var editRequest models.EditScheduleRequest
	if err := c.ShouldBindJSON(&editRequest); err != nil {
		log.Error("Error in EditScheduleHandler - ShouldBindJSON()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if editRequest.WorkDateString == "" && editRequest.ShiftLength <= 0 {
		log.Error("Error in EditScheduleHandler - Input Validation", errors.New(utils.ErrFieldMissing))
		utils.Fail(c, http.StatusBadRequest, utils.ErrFieldMissing)
		return
	}
	idStr := c.Param("schedule_id")
	if idStr == "" {
		log.Error("Error in EditScheduleHandler - schedule_id not present")
		utils.Fail(c, http.StatusBadRequest, utils.ErrInvalidSchedule)
		return
	}
	value, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("Error in EditScheduleHandler - parsing schedule_id", err.Error())
		utils.Fail(c, http.StatusBadRequest, utils.ErrBadInput)
		return
	}

	userID, err := utils.ExtractUserIDFromContextParam(c)
	if err != nil {
		log.Error("Error in EditScheduleHandler - ExtractUserIDFromContextParam()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	scheduleID := uint(value)
	schedule, err := h.Services.AdminService.EditSchedule(c, scheduleID, userID, &editRequest)
	if err != nil {
		log.Error("Error in EditScheduleHandler - service EditSchedule()", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Fail(c, http.StatusNotFound, utils.ErrRecordNotFound.Error())
			return
		}
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp := Response{
		Message: utils.MsgScheduleUpdated,
		Payload: schedule,
	}
	log.Info(utils.MsgScheduleUpdated)
	utils.Send(c, http.StatusOK, resp)
}

func (h *Handler) DeleteScheduleHandler(c *gin.Context) {
	idStr := c.Param("schedule_id")
	if idStr == "" {
		log.Error("Error in DeleteScheduleHandler - schedule_id not present")
		utils.Fail(c, http.StatusBadRequest, utils.ErrInvalidSchedule)
		return
	}

	value, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("Error in DeleteScheduleHandler - parsing schedule_id")
		utils.Fail(c, http.StatusBadRequest, utils.ErrBadInput)
		return
	}

	scheduleID := uint(value)
	err = h.Services.AdminService.DeleteSchedule(c, scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error in DeleteScheduleHandler - record not found", err.Error())
			utils.Fail(c, http.StatusNotFound, utils.ErrRecordNotFound.Error())
			return
		}
		log.Error("Error in DeleteScheduleHandler - service unknown error", err.Error())
		utils.Fail(c, http.StatusInternalServerError, err.Error())
	}
	resp := Response{
		Message: utils.MsgScheduleDeleted,
		Payload: gin.H{
			"schedule_id": scheduleID,
		},
	}
	log.Info(utils.MsgScheduleDeleted)
	utils.Send(c, http.StatusOK, resp)
}

func (h *Handler) GetUsersByAccumulatedHoursHandler(c *gin.Context) {
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	users, err := h.Services.AdminService.GetUsersByAccumulatedHours(c, start_date, end_date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error in GetUsersByAccumulatedHoursHandler - record not found", err.Error())
			utils.Fail(c, http.StatusNotFound, utils.ErrRecordNotFound.Error())
			return
		}
		log.Error("Error in GetUsersByAccumulatedHoursHandler - unknwon error", err.Error())
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp := Response{
		Message: utils.MsgGetUsersOrdered,
		Payload: users,
	}
	log.Info(utils.MsgGetUsersOrdered)
	utils.Send(c, http.StatusOK, resp)
}

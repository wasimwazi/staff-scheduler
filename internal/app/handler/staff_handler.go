package handler

import (
	"errors"
	"fmt"
	"net/http"
	"scheduler/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (h *Handler) GetScheduleHandler(c *gin.Context) {
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	staffIDStr := c.Query("staff_id")
	fmt.Println("hereee")

	var userID uint

	if staffIDStr != "" {
		staffID, err := strconv.ParseUint(staffIDStr, 10, 64)
		if err != nil {
			log.Error("Error in GetScheduleHandler - Parse staffID", err.Error())
			utils.Fail(c, http.StatusBadRequest, "Invalid staff_id provided")
			return
		}
		userID = uint(staffID)
	} else {
		userFromContext, exists := c.Get("currentUserID")
		if !exists {
			log.Error("Error in GetScheduleHandler - Get current user from context")
			utils.Fail(c, http.StatusBadRequest, "User not valid")
			return
		}
		currentUser, ok := userFromContext.(uint)
		if !ok {
			log.Error("Error in GetScheduleHandler - parsing the current user")
			utils.Fail(c, http.StatusBadRequest, utils.ErrParsingCurrentUser)
			return
		}
		userID = uint(currentUser)
	}

	scheduleResp, err := h.Services.StaffService.GetStaffSchedules(c, userID, start_date, end_date)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error in GetScheduleHandler - record not found")
			utils.Fail(c, http.StatusNotFound, utils.ErrRecordNotFound.Error())
			return
		}
		log.Error("Error in GetScheduleHandler - unknown error from service")
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp := Response{
		Message: "Staff schedule from " + start_date + " to " + end_date,
		Payload: scheduleResp,
	}
	log.Info(resp.Message)
	utils.Send(c, http.StatusOK, resp)
}

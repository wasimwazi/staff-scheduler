package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DateParser(startDateString, endDateString string) (*time.Time, *time.Time, error) {
	startDate, err := time.Parse(DateFormat, startDateString)
	if err != nil {
		return nil, nil, ErrInvalidDateInput
	}
	endDate, err := time.Parse(DateFormat, endDateString)
	if err != nil {
		return nil, nil, ErrInvalidDateInput
	}
	if startDate.IsZero() {
		sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
		startDate = sevenDaysAgo
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if startDate.Before(oneYearAgo) {
		return nil, nil, ErrStartDateOutOfRange
	}

	if startDate.After(endDate) {
		return nil, nil, ErrStartDateGreaterThanEndDate
	}

	return &startDate, &endDate, nil
}

func ExtractUserIDFromContextParam(c *gin.Context) (uint, error) {
	idStr := c.Param("user_id")
	if idStr == "" {
		return 0, errors.New("user_id field is empty")
	}

	value, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New(ErrBadInput)
	}
	userID := uint(value)
	return userID, nil
}

func ExtractCurrentUserIDFromContext(c *gin.Context) (uint, error) {
	userFromContext, exists := c.Get("currentUserID")
	if !exists {
		return 0, errors.New("user not valid")
	}

	currentUser, ok := userFromContext.(uint)
	if !ok {
		return 0, errors.New("error parsing the currentUserID from context")
	}
	return currentUser, nil
}

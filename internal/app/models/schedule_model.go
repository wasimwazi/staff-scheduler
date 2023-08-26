package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	WorkDate    time.Time
	User        string
	AccountID   uint
	ShiftLength int
}

type CreateScheduleRequest struct {
	UserID         uint      `json:"user_id,omitempty"`
	WorkDateString string    `json:"date" binding:"required"`
	ShiftLength    int       `json:"shift" binding:"required,min=1,max=12"`
	WorkDate       time.Time `json:"omitempty"`
}

type CreateScheduleResponse struct {
	UserID      uint      `json:"user_id" binding:"required"`
	ShiftLength int       `json:"shift" binding:"required,min=1,max=12"`
	WorkDate    time.Time `json:"date,omitempty"`
	ID          uint      `json:"schedule_id" binding:"required"`
	Username    string    `json:"username" binding:"required"`
}

type ScheduleModel struct {
	WorkDateString string    `json:"date" binding:"required"`
	ShiftLength    int       `json:"shift" binding:"required,min=1,max=12"`
	WorkDate       time.Time `json:"omitempty"`
}

type UserSchedules struct {
	UserID    uint            `json:"user_id" binding:"required"`
	Schedules []ScheduleModel `json:"schedules" binding:"required"`
}

type ScheduleRequest struct {
	Users []UserSchedules `json:"users" binding:"required"`
}

type EditScheduleRequest struct {
	WorkDateString string    `json:"date,omitempty"`
	ShiftLength    int       `json:"shift,omitempty"`
	WorkDate       time.Time `json:"omitempty"`
}

type ScheduleResponse struct {
	ID          uint   `json:"id,omitempty"`
	WorkDate    string `json:"date,omitempty"`
	ShiftLength int    `json:"shift_length,omitempty"`
	Name        string `json:"user,omitempty"`
}

type AccountWithAccumulatedShiftResponse struct {
	ID                   uint   `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Email                string `json:"email,omitempty"`
	Phone                string `json:"phone,omitempty"`
	AccumulatedWorkHours int    `json:"accumulated_work_hours,omitempty"`
}

type GetScheduleResponse struct {
	UserID    uint
	Schedules []UserScheduleResponse
}

type UserScheduleResponse struct {
	WorkDate    string `json:"work_date,omitempty"`
	ShiftLength int    `json:"shift_length,omitempty"`
}

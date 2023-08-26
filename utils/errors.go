package utils

import "errors"

const (
	DecodeError = "decode error, please check the input"

	ErrPasswordEmpty = "password is empty"

	ErrRoleRequired = "role is required"

	ErrPhoneValidation = "invalid phone number format (E.164)"

	ErrInvalidEmail = "invalid email input"

	ErrEmailValidation = "email is empty"

	ErrPassword = "wrong password"

	ErrJWT = "JWT Error"

	AccessTokenInterval = 100

	InvalidTokenError = "invalid Token"

	TopNError = "'number' parameter empty"

	ErrFieldMissing = "at least one field must be provided"

	ErrBadInput = "bad input"

	ErrCantDeletUser = "can't delete the current user"

	ErrInvalidSchedule = "invalid schedule id"

	ErrParsingCurrentUser = "Error parsing the currentUserID from context"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrPasswordCondition = errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")

	ErrUserInactive = errors.New("user is inactive")

	ErrInvalidDateInput = errors.New("invalid date input")

	ErrRecordNotFound = errors.New("record not found")

	ErrStartDateOutOfRange = errors.New("start date is out of range")

	ErrStartDateGreaterThanEndDate = errors.New("start date greater than end date")

	ErrScheduleAlreadyExist = errors.New("schedule already exist for the given dates")

	ErrCreateScheduleForOtherUsers = errors.New("can't create schedule for users other than staff")

	ErrFetchData = errors.New("error while fetching data")

	ErrValidatingInput = errors.New("error validating input")

	ErrInvalidRole = errors.New("invalid role input")
)

type ErrResponse struct {
	Error string `json:"error"`
}
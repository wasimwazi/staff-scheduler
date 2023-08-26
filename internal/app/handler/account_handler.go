package handler

import (
	"errors"
	"net/http"
	"scheduler/internal/app/models"
	"scheduler/utils"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// LoginAccount : to Login Account
func (h *Handler) LoginAccount(c *gin.Context) {
	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Error("Error in LoginAccount - ShouldBindJSON()", err.Error())
		utils.Fail(c, http.StatusBadRequest, utils.DecodeError)
		return
	}

	// Validate the login request
	if err := ValidateLoginRequest(login); err != nil {
		log.Error("Error in LoginAccount - ValidateLoginRequest()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	accountLogged, err := h.Services.AccountService.LoginAccount(c, &login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error in LoginAccount - record not found", err.Error())
			utils.Fail(c, http.StatusNotFound, utils.ErrRecordNotFound.Error())
			return
		} else if errors.Is(err, errors.New(utils.ErrPassword)) {
			log.Error("Error in LoginAccount - error password", err.Error())
			utils.Fail(c, http.StatusForbidden, utils.ErrPassword)
			return
		} else if errors.Is(err, utils.ErrUserInactive) {
			log.Error("Error in LoginAccount - user inactive", err.Error())
			utils.Fail(c, http.StatusForbidden, utils.ErrUserInactive.Error())
			return
		}
		log.Error("Error in LoginAccount - unknwon error from service", err.Error())
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	jwt, err := h.JWTGenerator.GenerateToken(accountLogged.ID, accountLogged.Email, accountLogged.Role)
	if err != nil {
		log.Error("Error in LoginAccount - GenerateToken()", err.Error())
		utils.Fail(c, http.StatusForbidden, err.Error())
		return
	}

	log.Info(utils.MsgLoggedIn)
	resp := Response{
		Message: utils.MsgLoggedIn,
		Payload: jwt,
	}
	utils.Send(c, http.StatusOK, resp)
}

//ValidateLogin : To validate the account logins
func ValidateLoginRequest(account models.LoginRequest) error {
	if len(account.Email) <= 0 {
		return errors.New(utils.ErrEmailValidation)
	}
	if len(account.Password) <= 0 {
		return errors.New(utils.ErrPasswordEmpty)
	}
	return nil
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var createRequest models.CreateAccountRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		var errMsg string
		if strings.Contains(err.Error(), "Email") {
			errMsg = utils.ErrInvalidEmail
		} else if strings.Contains(err.Error(), "Phone") && strings.Contains(err.Error(), "e164") {
			errMsg = utils.ErrPhoneValidation
		} else if strings.Contains(err.Error(), "Password") {
			errMsg = utils.ErrPasswordEmpty
		} else if strings.Contains(err.Error(), "Role") {
			errMsg = utils.ErrRoleRequired
		} else {
			errMsg = utils.ErrValidatingInput.Error()
		}

		log.Info("Error in create account", err)
		log.Error("Error in CreateAccount - validation error", err.Error())
		utils.Fail(c, http.StatusBadRequest, errMsg)
		return
	}

	if err := AccountRequestValidator(createRequest, createRequest.Role); err != nil {
		log.Error("Error in CreateAccount - AccountRequestValidator()", err.Error())
		utils.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.Services.AccountService.CreateAccount(c, &createRequest)
	if err != nil {
		if errors.Is(err, utils.ErrEmailAlreadyExists) {
			log.Error("Error in CreateAccount - email already exists()", err.Error())
			utils.Fail(c, http.StatusBadRequest, utils.ErrEmailAlreadyExists.Error())
			return
		}
		log.Error("Error in CreateAccount - unknown error from service", err.Error())
		utils.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := Response{
		Message: utils.MsgAccountCreated,
		Payload: account,
	}
	log.Info(utils.MsgAccountCreated)
	utils.Send(c, http.StatusCreated, resp)
}

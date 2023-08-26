package main

import (
	"scheduler/internal/app/handler"
	"scheduler/internal/app/service"
	"scheduler/internal/platform/database"
	"scheduler/utils"

	log "github.com/sirupsen/logrus"
)

// Run - sets up our application
func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting Up our APP")

	err := utils.CheckEnv()
	if err != nil {
		log.Fatal("Error while checking the environment files", err)
		panic(err)
	}

	store, err := database.NewDatabase()
	if err != nil {
		log.Error("failed to setup connection to the database")
		return err
	}
	err = database.DBMigration(store)
	if err != nil {
		log.Error("failed to setup database")
		return err
	}

	accountService := service.NewAccountService(store)
	adminService := service.NewAdminService(store)
	staffService := service.NewStaffService(store)

	services := service.NewServices(accountService, adminService, staffService)

	handler := handler.NewHandler(services)

	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}

	return nil
}


// @title Schedules API
func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}

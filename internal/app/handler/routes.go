package handler

import (
	"path/filepath"
	"scheduler/internal/app/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// mapRoutes - sets up all the routes for our application
func (h *Handler) mapRoutes() {
	h.Router.GET("/alive", h.AliveCheck)

	absPath, err := filepath.Abs("docs/")
	if err != nil {
		panic(err)
	}
	h.Router.Static("/docs", absPath)
	h.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/docs/swagger.yaml")))

	apiV1 := h.Router.Group("/api/v1")

	apiV1.POST("/login", h.LoginAccount)
	apiV1.POST("/account", h.CreateAccount)

	// Middleware to verify user token
	apiV1.Use(middleware.VerifyToken())

	h.mapStaffRoutes(apiV1)
	h.mapAdminRoutes(apiV1)
}

func (h *Handler) mapStaffRoutes(parent *gin.RouterGroup) {
	staffGroup := parent.Group("/staff")
	staffGroup.Use(middleware.RoleMiddleware("staff"))

	// get schedules
	staffGroup.GET("/schedule", h.GetScheduleHandler)
}

func (h *Handler) mapAdminRoutes(parent *gin.RouterGroup) {
	adminGroup := parent.Group("/admin")
	adminGroup.Use(middleware.RoleMiddleware("admin"))

	// User management
	adminGroup.PUT("/users/:user_id", h.EditUserHandler)
	adminGroup.DELETE("/users/:user_id", h.DeleteUserHandler)

	// Schedule use cases
	adminGroup.POST("/users/:user_id/schedule", h.CreateScheduleHandler)
	adminGroup.PUT("users/:user_id/schedule/:schedule_id", h.EditScheduleHandler)
	adminGroup.DELETE("users/:user_id/schedule/:schedule_id", h.DeleteScheduleHandler)
	adminGroup.GET("/users", h.GetUsersByAccumulatedHoursHandler)
}

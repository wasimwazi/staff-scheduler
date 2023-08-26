package handler

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"scheduler/internal/app/middleware"
	"scheduler/internal/app/service"
	"scheduler/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router       *gin.Engine
	Services     *service.Services
	Server       *http.Server
	JWTGenerator utils.TokenGenerator
}

// Response objecgi
type Response struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload,omitempty"`
}

// NewHandler - returns a pointer to a Handler
func NewHandler(services *service.Services) *Handler {
	log.Info("setting up our handler")

	jwtGenerator := utils.NewJWTGenerator()
	h := &Handler{
		Services:     services,
		Router:       gin.Default(),
		JWTGenerator: jwtGenerator,
	}

	h.Router.Use(middleware.JSONMiddleware())
	h.Router.Use(middleware.LoggingMiddleware())
	h.Router.Use(middleware.TimeoutMiddleware())

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.Router,
	}
	return h
}

func (h *Handler) AliveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am Alive!"})
}

// Serve - gracefully serves our newly set up handler function
func (h *Handler) Serve() error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Initiate server shutdown and wait for it to complete
	if err := h.Server.Shutdown(ctx); err != nil {
		log.Println("Error during server shutdown:", err)
	}

	wg.Wait()

	log.Println("Server has shut down gracefully")
	return nil
}

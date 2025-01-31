package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fox1N69/iq-testtask/internal/config"
	"github.com/Fox1N69/iq-testtask/internal/delivery/http/handler"
	"github.com/Fox1N69/iq-testtask/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config          *config.Config
	router          *gin.Engine
	userHandler     handler.UserHandler
	transactHandler handler.TransactionHandler
	server          *http.Server
	log             logger.Logger
}

func New(
	config *config.Config,
	userHandler handler.UserHandler,
	transactHandler handler.TransactionHandler,
) (*Server, error) {
	router := gin.Default()
	server := &Server{
		config:          config,
		router:          router,
		userHandler:     userHandler,
		transactHandler: transactHandler,
		server: &http.Server{
			Handler: router,
		},
		log: logger.GetLogger(),
	}

	server.routes()

	return server, nil
}

func (s *Server) routes() {
	api := s.router.Group("/api")
	{

		user := api.Group("/user")
		{
			user.GET("/:id", s.userHandler.UserByID)
		}

		transact := api.Group("/transaction")
		{
			transact.GET("/:user_id", s.transactHandler.LastTransactions)
			transact.POST("/transfer", s.transactHandler.Transfer)
			transact.POST("/replenish", s.transactHandler.Replenish)
		}
	}
}

func (s *Server) Start() error {
	s.log.Info("Server start on port:", s.config.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			s.log.Errorf("listen %s\n", err)
		}
	}()

	<-quit
	s.log.Info("Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Fatal("Server Shutdown:", err)
	}
	s.log.Info("Server exiting")

	return nil
}

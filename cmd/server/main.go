package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/auth"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/config"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/healthcheck"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/task"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfgPath := "config/local.yml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		cfgPath = envPath
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		fmt.Printf("failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("failed to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to database")

	// Setup layers
	authRepo := repository.NewAuthRepository(db)
	authService := auth.NewService(authRepo, cfg.JWT.Secret)

	taskRepo := repository.NewTaskRepository(db)
	taskService := task.NewService(taskRepo)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register handlers
	healthcheck.RegisterHandlers(r)
	auth.RegisterHandlers(r, authService)
	task.RegisterHandlers(r, taskService)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("server error: %v\n", err)
		os.Exit(1)
	}
}

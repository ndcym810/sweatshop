// internal/app/app.go
package app

import (
	"github.com/sweatshop/sweatshop/internal/claudeteam"
	"github.com/sweatshop/sweatshop/internal/department"
	"github.com/sweatshop/sweatshop/internal/project"
	"github.com/sweatshop/sweatshop/internal/task"
	"github.com/sweatshop/sweatshop/internal/team"
	"github.com/sweatshop/sweatshop/pkg/logger"
)

// App holds application dependencies
type App struct {
	TeamHandler       *team.Handler
	ProjectHandler    *project.Handler
	DepartmentHandler *department.Handler
	TaskHandler       *task.Handler
	ClaudeTeamHandler *claudeteam.Handler
	WebSocketHub      *claudeteam.WebSocketHub
	Watcher           *claudeteam.Watcher
}

// New creates a new App instance
func New() *App {
	// Initialize services
	teamSvc := team.NewService()
	projectSvc := project.NewService()
	departmentSvc := department.NewService()
	taskSvc := task.NewService()
	claudeTeamSvc := claudeteam.NewService()

	// Initialize WebSocket hub
	wsHub := claudeteam.NewWebSocketHub()
	go wsHub.Run()

	// Initialize file watcher
	watcher, err := claudeteam.NewWatcher(claudeTeamSvc, wsHub)
	if err != nil {
		logger.Error.Printf("Failed to create watcher: %v", err)
	} else {
		if err := watcher.Start(); err != nil {
			logger.Error.Printf("Failed to start watcher: %v", err)
		}
	}

	// Initialize handlers
	app := &App{
		TeamHandler:       team.NewHandler(teamSvc),
		ProjectHandler:    project.NewHandler(projectSvc),
		DepartmentHandler: department.NewHandler(departmentSvc),
		TaskHandler:       task.NewHandler(taskSvc),
		ClaudeTeamHandler: claudeteam.NewHandler(claudeTeamSvc),
		WebSocketHub:      wsHub,
		Watcher:           watcher,
	}

	logger.Info.Println("Application initialized")
	return app
}

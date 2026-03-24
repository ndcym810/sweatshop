// internal/app/app.go
package app

import (
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
}

// New creates a new App instance
func New() *App {
	// Initialize services
	teamSvc := team.NewService()
	projectSvc := project.NewService()
	departmentSvc := department.NewService()
	taskSvc := task.NewService()

	// Initialize handlers
	app := &App{
		TeamHandler:       team.NewHandler(teamSvc),
		ProjectHandler:    project.NewHandler(projectSvc),
		DepartmentHandler: department.NewHandler(departmentSvc),
		TaskHandler:       task.NewHandler(taskSvc),
	}

	logger.Info.Println("Application initialized")
	return app
}

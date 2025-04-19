package app

import (
	"zhangyumao/config"
	"zhangyumao/internal/github"

	"ergo.services/ergo/gen"
)

func CreateApp() gen.ApplicationBehavior {
	return &App{}
}

type App struct{
	config config.GitHubAppConfig
}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *App) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	app.config = config.GetGitHubAppConfig()
	clientCreator := github.NewClientCreator(app.config)

	return gen.ApplicationSpec{
		Name:        "app",
		Description: "description of this application",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    alias_Handler,
				Factory: factory_Handler,
				Args: []any{
					HandlerInitArgs{ClientCreator: &clientCreator},
				},
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *App) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *App) Terminate(reason error) {}

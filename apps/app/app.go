package app

import (
	"ergo.services/ergo/gen"
)

func CreateApp() gen.ApplicationBehavior {
	return &App{}
}

type App struct{}

// Load invoked on loading application using method ApplicationLoad of gen.Node interface.
func (app *App) Load(node gen.Node, args ...any) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "app",
		Description: "description of this application",
		Mode:        gen.ApplicationModeTransient,
		Group: []gen.ApplicationMemberSpec{
			{
				Name:    alias_AppsClient,
				Factory: factory_AppsClient,
			},
			{
				Name:    alias_Handler,
				Factory: factory_Handler,
			},
		},
	}, nil
}

// Start invoked once the application started
func (app *App) Start(mode gen.ApplicationMode) {}

// Terminate invoked once the application stopped
func (app *App) Terminate(reason error) {}

package app

type ServInterface interface {
	Init(app *App) error
	Run(app *App) error
	Shutdown(app *App) error
	Reload(app *App) error
}

package internalhttp

type Handler struct {
	app    Application
	logger Logger
}

func NewHandler(app Application, logger Logger) *Handler {
	return &Handler{
		app:    app,
		logger: logger,
	}
}

package health

import (
	"net/http"

	"github.com/asciifaceman/godock/pkg/model"

	"github.com/asciifaceman/godock/pkg/framework"
	"github.com/asciifaceman/godock/pkg/handlers"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var _ handlers.Handler = (*Handler)(nil)

// Handler ...
type Handler struct {
	handlers.BaseHandler
	Logger *zap.Logger
}

// NewHealthHandler ...
func NewHealthHandler(log *zap.Logger) *Handler {
	return &Handler{
		Logger: log,
	}
}

// Routes returns a list of routes for the server to bind
func (h *Handler) Routes() []framework.IRoute {
	return []framework.IRoute{
		&framework.Route{
			Path:     "/health",
			Method:   framework.GET.String(),
			Request:  nil,
			Response: &model.HealthCheckResponse{},
			Handler:  h.GetHealth,
		},
	}
}

// GetHealth ...
func (h *Handler) GetHealth(e echo.Context) error {
	return h.SendResponse(e, http.StatusOK, &model.HealthCheckResponse{
		Status: "ok",
	})
}

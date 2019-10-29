package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/asciifaceman/godock/pkg/framework"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Handler is an interface describing request handlers
type Handler interface {
	Routes() []framework.IRoute
}

var _ Handler = (*BaseHandler)(nil)

// BaseHandlerError is a basic handler error response model
type BaseHandlerError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// BaseHandler ...
type BaseHandler struct {
	logger     *zap.Logger
	Registered bool
}

// Routes returns a list of routes for the server to bind
func (b *BaseHandler) Routes() []framework.IRoute {
	return []framework.IRoute{}
}

// JSONResponse writes a JSON response
func (b BaseHandler) JSONResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// SendResponse wraps handler responses into a jsend format
func (b *BaseHandler) SendResponse(c echo.Context, code int, i interface{}) error {
	return c.JSON(code, i)
}

// SendError wraps error responses into a jsend format
func (b *BaseHandler) SendError(c echo.Context, code int, e error) error {
	return c.JSON(code, &BaseHandlerError{
		Status:  "error",
		Message: e.Error(),
	})
}

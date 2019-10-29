package framework

import (
	echo "github.com/labstack/echo/v4"
)

var _ IRoute = (*Route)(nil)

// IRoute ...
type IRoute interface {
	GetPath() string
	GetMethod() string
	GetHandler() echo.HandlerFunc
}

// Route is the base structure of an API route
type Route struct {
	Path    string
	Handler echo.HandlerFunc

	Method string

	Request  interface{}
	Response interface{}
}

// GetPath returns a route's path
func (r *Route) GetPath() string {
	return r.Path
}

// GetMethod returns a route's method
func (r *Route) GetMethod() string {
	return r.Method
}

// GetHandler returns a route's handler
func (r *Route) GetHandler() echo.HandlerFunc {
	return r.Handler
}

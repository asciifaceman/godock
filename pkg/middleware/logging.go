package middleware

import (
	"time"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggingMiddleware ...
func LoggingMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err = next(c); err != nil {
				c.Error(err)
			}

			resID := res.Header().Get(echo.HeaderXRequestID)

			stop := time.Now()
			lat := float64(stop.Sub(start)) / float64(time.Millisecond)
			m := req.Method
			s := res.Status
			ip := req.RemoteAddr
			ua := req.UserAgent()
			p := req.URL.Path

			logger.Info("Request processed",
				zap.String("ID", resID),
				zap.Float64("Duration", lat),
				zap.String("Method", m),
				zap.Int("Status", s),
				zap.String("Path", p),
				zap.String("Requester", ip),
				zap.String("User-Agent", ua),
			)

			return nil
		}
	}
}

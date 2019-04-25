package rover

import (
	"github.com/labstack/echo"
	"time"
)

// SendTimingMetrics : An Echo middleware to add timing metrics around each Echo Request.
func (rover *Rover) SendTimingMetrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) (err error) {
			requestStart := time.Now()
			err = next(context)
			duration := time.Since(requestStart)

			requestKey := rover.generateRequestKey(context)
			tags := rover.getTagsForRequest(context)
			rover.StatsdClient.Timing(requestKey, duration, tags, 1)

			return
		}
	}
}

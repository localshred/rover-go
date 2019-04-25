package rover

import (
	"github.com/labstack/echo"
)

const (
	statsdRequestKeyContextKey = "rover.statsd.requestKey"
)

// RequestKeyMiddleware : An Echo middleware to add timing metrics around each Echo Request.
func (rover *Rover) RequestKeyMiddleware(requestKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) (err error) {
			setRequestKeyOnContext(context, requestKey)
			return next(context)
		}
	}
}

func getRequestKeyFromContext(context echo.Context) interface{} {
	return context.Get(statsdRequestKeyContextKey)
}

func setRequestKeyOnContext(context echo.Context, requestKey string) {
	context.Set(statsdRequestKeyContextKey, requestKey)
}

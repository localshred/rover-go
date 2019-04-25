package rover

import (
	"github.com/labstack/echo"
)

const (
	statsdRequestKeyContextKey = "rover.statsd.requestKey"
)

// SetRequestKey : An Echo middleware to set the static requestKey on the echo
// context for later retrieval by other middleware (see TimingMiddleware).
func (rover *Rover) SetRequestKey(requestKey string) echo.MiddlewareFunc {
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

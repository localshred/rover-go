package rover

import (
	"github.com/labstack/echo"
)

const (
	statsdRequestPageContextKey = "rover.statsd.requestPage"
)

// SetRequestPage : An Echo middleware to set the static requestPage on the echo
// context for later retrieval by other middleware (see TimingMiddleware).
func (rover *Rover) SetRequestPage(requestPage string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) (err error) {
			setRequestPageOnContext(context, requestPage)
			return next(context)
		}
	}
}

func getRequestPageFromContext(context echo.Context) interface{} {
	return context.Get(statsdRequestPageContextKey)
}

func setRequestPageOnContext(context echo.Context, requestPage string) {
	context.Set(statsdRequestPageContextKey, requestPage)
}

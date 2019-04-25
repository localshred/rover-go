package rover

import (
	"errors"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo"
)

type (
	// Rover : Rover config
	Rover struct {
		StatsdClient *statsd.Client

		// AddTagsForRequest is a callback you can provide to generate your own tags for each request to be
		// sent with the timing metrics. The tags will be requested after the main request handler has
		// completed its work.
		AddTagsForRequest func(context echo.Context) []string

		// GetRequestKey is a callback that returns the request key to be appended to the statsd client's
		// Namespace (if any). If this callback is not provided, the request key will be extracted from the
		// echo context if available. Otherwise, a request key will be auto-generated based on the Request URI.
		//
		// The preferred way of specifying the request key is to use the config.RequestKeyMiddleware function
		// to statically set the rqeuest key on the echo context. See that function for more info.
		//
		// If the Namespace was "mycompany.myapp.prod" and GetRequestKey returns "get.list_users.timing", the full
		// keyspace will be "mycompany.myapp.prod.get./users.timing".
		GetRequestKey func(context echo.Context) string
	}
)

// New : Creates a new rover configuration with the given statsd client and sensible defaults.
func New(statsdClient *statsd.Client) *Rover {
	return &Rover{
		StatsdClient:  statsdClient,
		GetRequestKey: defaultGetRequestKey,
	}
}

// Tag : Builds a tag given a namespace and value and joining them with a colon.
//
//       Tag("foo:bar", "baz") -> "foo:bar:baz"
//
func Tag(namespace, value string) string {
	return fmt.Sprintf("%s:%s", namespace, value)
}

func defaultGetRequestKey(context echo.Context) string {
	if requestKey := getRequestKeyFromContext(context); requestKey != nil {
		return requestKey.(string)
	}
	return context.Request().URL.Path
}

func (rover *Rover) generateRequestKey(context echo.Context) string {
	if rover.GetRequestKey == nil {
		panic(errors.New("GetRequestKey callback cannot be nil"))
	}
	return rover.GetRequestKey(context)
}

func (rover *Rover) getTagsForRequest(context echo.Context) []string {
	request := context.Request()
	tags := []string{
		Tag("http:host", request.Host),
		Tag("http:method", request.Method),
		Tag("http:remote-addr", request.RemoteAddr),
	}
	if rover.AddTagsForRequest != nil {
		return append(tags, rover.AddTagsForRequest(context)...)
	}
	return tags
}

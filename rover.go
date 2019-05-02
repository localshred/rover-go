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

		// GetRequestPage is a callback that returns the request key to be appended to the statsd client's
		// Namespace (if any). If this callback is not provided, the request key will be extracted from the
		// echo context if available. Otherwise, a request key will be auto-generated based on the Request URI.
		//
		// The preferred way of specifying the request key is to use the config.RequestPageMiddleware function
		// to statically set the rqeuest key on the echo context. See that function for more info.
		//
		// If the Namespace was "mycompany.myapp.prod" and GetRequestPage returns "get.list_users.timing", the full
		// keyspace will be "mycompany.myapp.prod.get./users.timing".
		GetRequestPage func(context echo.Context) string
	}
)

// New : Creates a new rover configuration with the given statsd client and sensible defaults.
func New(statsdClient *statsd.Client) *Rover {
	return &Rover{
		StatsdClient:   statsdClient,
		GetRequestPage: defaultGetRequestPage,
	}
}

// Tag : Builds a tag given a namespace and value and joining them with a colon.
//
//       Tag("foo:bar", "baz") -> "foo:bar:baz"
//
func Tag(namespace, value string) string {
	return fmt.Sprintf("%s:%s", namespace, value)
}

func defaultGetRequestPage(context echo.Context) string {
	if requestPage := getRequestPageFromContext(context); requestPage != nil {
		return requestPage.(string)
	}
	return context.Request().URL.Path
}

func (rover *Rover) generateRequestPage(context echo.Context) string {
	if rover.GetRequestPage == nil {
		panic(errors.New("GetRequestPage callback cannot be nil"))
	}
	return rover.GetRequestPage(context)
}

func (rover *Rover) getTagsForRequest(context echo.Context) []string {
	request := context.Request()
	tags := []string{
		Tag("http:host", request.Host),
		Tag("http:method", request.Method),
	}
	if rover.AddTagsForRequest != nil {
		return append(tags, rover.AddTagsForRequest(context)...)
	}
	return tags
}

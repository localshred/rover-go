# Rover

Provides middleware and utilities for integrating DataDog statsd with the Echo
HTTP library written for go servers.

### SendTimingMetrics

The `SendTimingMetrics` is used for any and all routes that your server will serve.
It simply times each route it wraps and sends a `timing` metric off to statsd
for the given route key.

You can specify which route key on a per-route basis via the `GetRequestKey` callback, usage of
`rover.SetRequestKey`, or letting rover utilize the default behavior (using the
request Path as the key).

```go
import (
  "github.com/DataDog/datadog-go/statsd"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/localshred/rover-go"
)

statsdClient, err := statsd.New("127.0.0.1:8125")
roverConfig := rover.New(statsdClient)

rover.GetRequestKey = func(context echo.Context) string {
  return context.Request().URL.Path // this is currently the default behavior
}

server := echo.New()
server.use(roverConfig.SendTimingMetrics())
server.use(middleware.CORS())
// ...
```

### SetRequestKey

The `SetRequestKey` is used to provide a static Request Key for each
route. The value is stored on the echo context and is retrieved by rover and
used if present. The middleware alone doesn't send anything to statsd, thus it's
really only useful in conjunction with the `SendTimingMetrics`.

```go
import (
  "github.com/DataDog/datadog-go/statsd"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/localshred/rover-go"
)

statsdClient, err := statsd.New("127.0.0.1:8125")
roverConfig := rover.New(statsdClient)

roverConfig.GetRequestKey = func(context echo.Context) string {
  return context.Request().URL.Path // this is currently the default behavior
}

server := echo.New()
server.use(roverConfig.SendTimingMetrics())
server.use(middleware.CORS())

server.POST("/user/login", roverConfig.SetRequestKey("user_login"), userLoginHandler)
server.POST("/user/logout", roverConfig.SetRequestKey("user_logout"), userLogoutHandler)
server.GET("/some/long/resource/url", roverConfig.SetRequestKey("foobar"), fooBarHandler)
// ...
```

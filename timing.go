package rover

import (
	"github.com/labstack/echo"
	"time"
)

const (
	defaultMetricName = "request_timing"
)

type (
	// TimingConfig : provides configuration to the timing middleware.
	TimingConfig struct {
		metricName string
	}
)

// SendTimingMetrics : An Echo middleware to add timing metrics around each Echo Request.
func (rover *Rover) SendTimingMetrics(config *TimingConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) (err error) {
			requestStart := time.Now()
			err = next(context)
			duration := time.Since(requestStart)

			metricName := getMetricName(config)
			requestPage := rover.generateRequestPage(context)
			tags := rover.getTagsForRequest(context)
			tags = append([]string{
				Tag("page", requestPage),
			})
			rover.StatsdClient.Timing(metricName, duration, tags, 1)

			return
		}
	}
}

func getMetricName(config *TimingConfig) string {
	if config != nil && config.metricName != "" {
		return config.metricName
	}
	return defaultMetricName
}

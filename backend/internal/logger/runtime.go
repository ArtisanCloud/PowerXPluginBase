package logger

import "github.com/sirupsen/logrus"

// WithRuntimeFields enriches the log entry with standard runtime metadata.
func WithRuntimeFields(pluginID, tenantID, traceID, component string, extra Fields) *logrus.Entry {
	fields := Fields{
		"plugin_id": pluginID,
		"tenant_id": tenantID,
		"trace_id":  traceID,
		"component": component,
	}
	for k, v := range extra {
		fields[k] = v
	}
	return WithFields(fields)
}

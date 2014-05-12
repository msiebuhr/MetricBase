package frontends

import (
	"github.com/msiebuhr/MetricBase/backends"
)

type Frontend interface {
	// SetBackend sets the server's backend
	SetBackend(backends.Backend)
	// Start the frontend (usually in a go-routine)
	Start()
	// Stop a previously started frontend
	Stop()
}

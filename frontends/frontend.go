package frontends

import (
	"github.com/msiebuhr/MetricBase/backends"
)

type Frontend interface {
	SetBackend(backends.Backend)
	Start()
	Stop()
}

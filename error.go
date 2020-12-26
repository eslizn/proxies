package proxies

import "errors"

var (
	ErrProbeFail        = errors.New("probe fail")
	ErrTransportsClosed = errors.New("transports already closed")
	ErrAssignTimeout    = errors.New("assign timeout")
	ErrTransportInvalid = errors.New("transport invalid")
)

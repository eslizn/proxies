package proxies

import (
	"sync/atomic"
	"time"
)

type Stats struct {
	Success  uint32
	Failure  uint32
	Verified *time.Time
}

func (s *Stats) Report(err error) {
	if err != nil {
		atomic.AddUint32(&s.Failure, 1)
	} else {
		atomic.AddUint32(&s.Success, 1)
	}
}

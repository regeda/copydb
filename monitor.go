package copydb

import "log"

// Monitor watches on events and errors.
type Monitor interface {
	VersionConflictDetected(error)
	ApplyFailed(error)
	PurgeFailed(error)
}

var defaultMonitor = monitor{}

type monitor struct{}

func (monitor) VersionConflictDetected(err error) {
	log.Printf("ERROR: %v", err)
}

func (monitor) ApplyFailed(err error) {
	log.Printf("ERROR: %v", err)
}

func (monitor) PurgeFailed(err error) {
	log.Printf("ERROR: %v", err)
}

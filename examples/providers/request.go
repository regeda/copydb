package providers

import "time"

// Request represents provider request update.
type Request struct {
	ID       string
	Remove   bool
	Set      map[string][]byte
	Unset    []string
	Currtime time.Time
}

package drones

import "time"

// Request represents drone request update.
type Request struct {
	ID       string
	Remove   bool
	Set      map[string][]byte
	Unset    []string
	Currtime time.Time
}

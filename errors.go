package copydb

import "errors"

var (
	ErrVersionConflict = errors.New("version conflict")
	ErrVersionNotFound = errors.New("version not found")
)

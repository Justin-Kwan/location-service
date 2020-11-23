package errors

import "github.com/pkg/errors"

var (
	// Drain errors
	DrainOutputNil = errors.New("Drain ouput channel is nil")

	// Redis DB errors
	KeyNotFound    = errors.New("Key passed in does not exist")
	MemberNotFound = errors.New("Member passed in does not exist")
	NoBatchQueries = errors.New("No geo queries passed into batch operation")

	// implement this!!
	// type error interface {
	// 	Error() string
	// }

	InvalidSearchRadius = errors.New("Invalid search radius given")
)

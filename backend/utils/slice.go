package utils

// EnsureNotNil returns an initialised empty slice when s is nil.
// This prevents JSON marshalling from encoding nil slices as null.
func EnsureNotNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

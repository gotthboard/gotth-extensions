package extensions

import "errors"

var (
	ErrInvalidInput        = errors.New("extensions: invalid input")
	ErrUnsupportedSchema   = errors.New("extensions: unsupported schema")
	ErrIncompatibleVersion = errors.New("extensions: incompatible version")
	ErrStaleBinding        = errors.New("extensions: stale binding")
	ErrUngrantedCapability = errors.New("extensions: ungranted capability")
	ErrUngrantedInterface  = errors.New("extensions: ungranted interface")
	ErrUngrantedSecret     = errors.New("extensions: ungranted secret")
	ErrInvalidTransition   = errors.New("extensions: invalid transition")
)

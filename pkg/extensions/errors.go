package extensions

import "errors"

var (
	// ErrInvalidInput classifies malformed, duplicate, or out-of-bound input.
	ErrInvalidInput = errors.New("extensions: invalid input")
	// ErrUnsupportedSchema classifies a manifest or grant schema outside V1.
	ErrUnsupportedSchema = errors.New("extensions: unsupported schema")
	// ErrIncompatibleVersion classifies a missing compatible control version.
	ErrIncompatibleVersion = errors.New("extensions: incompatible version")
	// ErrStaleBinding classifies a grant that does not bind the supplied manifest.
	ErrStaleBinding = errors.New("extensions: stale binding")
	// ErrUngrantedCapability classifies a capability outside the manifest
	// declaration.
	ErrUngrantedCapability = errors.New("extensions: ungranted capability")
	// ErrUngrantedInterface classifies an interface outside declared or supported
	// ranges.
	ErrUngrantedInterface = errors.New("extensions: ungranted interface")
	// ErrUngrantedSecret classifies an unknown, absent, or otherwise ungranted
	// secret slot.
	ErrUngrantedSecret = errors.New("extensions: ungranted secret")
	// ErrInvalidTransition classifies an unknown state or forbidden lifecycle edge.
	ErrInvalidTransition = errors.New("extensions: invalid transition")
)

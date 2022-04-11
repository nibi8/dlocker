package models

import (
	"errors"
)

var ErrNotFound = errors.New("not_found")
var ErrDuplicate = errors.New("duplicate")
var ErrNoLuck = errors.New("no_luck")
var ErrNotSupported = errors.New("not_supported")

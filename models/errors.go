package models

import (
	"errors"
)

var ErrNotFound = errors.New("not_found")
var ErrDuplicate = errors.New("duplicate")
var ErrNoLuck = errors.New("no_luck")
var ErrNotSupported = errors.New("not_supported")
var ErrBadParam = errors.New("bad_param")
var ErrBadContext = errors.New("bad_context")
var ErrWrongState = errors.New("wrong_state")

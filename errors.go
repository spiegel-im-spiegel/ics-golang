package ics

import "errors"

var (
	ErrNoEvent    = errors.New("no event")
	ErrNoEventDay = errors.New("no event for the day")
)

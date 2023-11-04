package model

import "errors"

var (
	ErrOriginConflict       = errors.New("origin url already exists")
	ErrShortConflict        = errors.New("short url already exists")
	ErrShortBadParam        = errors.New("short url doesn't exist")
	ErrUrlGeneratorInternal = errors.New("short url generator closed channel")
)

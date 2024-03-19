package storage

import (
	"errors"
)

var (
	ErrNotUniqueAction = errors.New("user already made such operation")
	ErrNotUniqueQuest  = errors.New("such quest already exists in the database")
)

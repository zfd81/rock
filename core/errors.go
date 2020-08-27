package core

import "github.com/pkg/errors"

var (
	ErrServExists   = errors.New("The service already exists")
	ErrServNotExist = errors.New("Service does not exist")
)

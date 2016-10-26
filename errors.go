package gorpc

import (
	"errors"
)

var (
	ErrPasswordIncorrect = errors.New("Password incorrect!")
	ErrURLInvalid = errors.New("URL error!")
)
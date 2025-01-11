package common

import (
	"errors"
)

var (
	ErrorCorrupt  = errors.New("Disturbance in the Void")
	ErrorFull     = errors.New("Void is full, extend?")
	ErrorInvalid  = errors.New("Form is not Void")
	ErrorNotFound = errors.New("Not found")
)

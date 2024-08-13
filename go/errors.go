package verbeux

import (
	"errors"
)

var (
	ErrCreateSession = errors.New("failed to create session")
	ErrSendMessage   = errors.New("failed to send message")
)

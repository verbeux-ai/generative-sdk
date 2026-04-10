package verbeux

import (
	"errors"
)

var (
	ErrCreateSession = errors.New("failed to create session")
	ErrReadSession   = errors.New("failed to read session")
	ErrSendMessage   = errors.New("failed to send message")
	ErrDeleteMessage = errors.New("failed to delete message")
)

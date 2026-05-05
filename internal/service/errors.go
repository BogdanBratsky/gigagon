package service

import "errors"

var (
	ErrRoomNotFound    = errors.New("room does not exist")
	ErrAlreadyStarted  = errors.New("room is already started")
	ErrNotHost         = errors.New("player is not host")
	ErrInvalidState    = errors.New("state is invalid")
	ErrAlreadyAnswered = errors.New("already answered")
	ErrAlreadyVoted    = errors.New("already voted")
	ErrInvalidAnswer   = errors.New("answer is invalid")
	ErrCantVoteSelf    = errors.New("you cannot vote yourself")
)

package jtt

import (
	"errors"
)

var (
	// ErrBodyTooLong too long message body
	ErrBodyTooLong = errors.New("too long message body")
	// ErrInvalidBody invalid message body
	ErrInvalidBody = errors.New("invalid message body")
	// ErrInvalidHeader invalid message header
	ErrInvalidHeader = errors.New("invalid message header")
	// ErrInvalidMessage invalid message format
	ErrInvalidMessage = errors.New("invalid message format")
	// ErrInvalidCheckSum invalid message check sum
	ErrInvalidCheckSum = errors.New("invalid message check sum")
	// ErrMethodNotImplemented method not implemented
	ErrMethodNotImplemented = errors.New("method not implemented")
	// ErrMessageNotRegistered message not registered
	ErrMessageNotRegistered = errors.New("message not registered")
	// ErrEntityDecode entity decode error
	ErrEntityDecode = errors.New("entity decode error")
	// ErrInvalidExtraLength invalid extra length
	ErrInvalidExtraLength = errors.New("invalid extra length")
	// ErrSegmentNotCompleted segment not completed
	ErrSegmentNotCompleted = errors.New("segment not completed")
)

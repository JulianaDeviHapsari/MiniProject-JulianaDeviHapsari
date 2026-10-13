package domain

import "errors"

// Sentinel kinds untuk error domain. Layer delivery memetakannya ke status HTTP.
var (
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
	ErrBadRequest   = errors.New("bad request")
)

// DomainError membawa pesan yang aman ditampilkan ke client.
type DomainError struct {
	Kind error
	Msg  string
}

func (e *DomainError) Error() string        { return e.Msg }
func (e *DomainError) Is(target error) bool { return target == e.Kind }

func NotFound(msg string) error     { return &DomainError{Kind: ErrNotFound, Msg: msg} }
func Forbidden(msg string) error    { return &DomainError{Kind: ErrForbidden, Msg: msg} }
func Unauthorized(msg string) error { return &DomainError{Kind: ErrUnauthorized, Msg: msg} }
func Conflict(msg string) error     { return &DomainError{Kind: ErrConflict, Msg: msg} }
func BadRequest(msg string) error   { return &DomainError{Kind: ErrBadRequest, Msg: msg} }

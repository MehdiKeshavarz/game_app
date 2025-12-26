package richerror

import "errors"

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type Op string
type RichError struct {
	operation    Op
	wrappedError error
	message      string
	kind         Kind
	meta         map[string]interface{}
}

func New(op Op) RichError {
	return RichError{operation: op}
}

func (r RichError) SetMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) SetKind(kind Kind) RichError {
	r.kind = kind
	return r
}

func (r RichError) SetWrappedError(wrappedError error) RichError {
	r.wrappedError = wrappedError
	return r
}

func (r RichError) SetMeta(meta map[string]interface{}) RichError {
	r.meta = meta
	return r
}

func (r RichError) Error() string {
	return r.message
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	var re RichError
	ok := errors.As(r.wrappedError, &re)
	if !ok {
		return 0
	}
	return re.Kind()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	var re RichError
	ok := errors.As(r.wrappedError, &re)
	if !ok {
		return ""
	}
	return re.Message()
}

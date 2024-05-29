package errorsx

import "errors"

type Errorx struct {
	code       int
	isInternal bool
	message    string
}

func (e *Errorx) Error() string {
	return e.message
}

func (e *Errorx) Internal() bool {
	return e.isInternal
}

func (e *Errorx) Code() int {
	return e.code
}

func New(isInternal bool, code int, message string) *Errorx {
	return &Errorx{
		isInternal: isInternal,
		code:       code,
		message:    message,
	}
}

func As(err error) (*Errorx, bool) {
	var rErr *Errorx
	ok := errors.As(err, &rErr)
	return rErr, ok
}

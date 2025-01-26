package exceptions

import "fmt"

type FileNotFound struct {
	origin_err error
	msg        string
}

func (e *FileNotFound) Error() string {
	return fmt.Sprintf("file not found. Origin err=%v, err_type=%T, msg=%s", e.origin_err, e.origin_err, e.msg)
}

type Option func(e *FileNotFound)

func WithOriginError(origin_err error) Option {
	return func(e *FileNotFound) { e.origin_err = origin_err }
}

func WithMsg(msg string) Option {
	return func(e *FileNotFound) { e.msg = msg }
}

func NewFileNotFound(opts ...Option) error {
	e := &FileNotFound{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

package service

import (
	"errors"
	"fmt"
)

type Err struct {
	SvcName string
	Method  string
	Type    ErrType
	err     error
}

type ErrType string

const (
	ErrTypeRequestValidation ErrType = "ERR_TYPE_REQUEST_VALIDATION"
	ErrTypeStorage           ErrType = "ERR_TYPE_STORAGE"
	ErrTypeBusinessLogic     ErrType = "ERR_TYPE_BUSINESS_LOGIC"
)

func (e *Err) Error() string {
	if e.err == nil {
		e.err = errNil
	}

	return fmt.Sprintf("error in %s.%s of type [%s]: %s", e.SvcName, e.Method, e.Type, e.err.Error())
}

func (e *Err) Unwrap() error {
	return e.err
}

func WrapError(e *Err, err error) error {
	e.err = err
	return e
}

// func (e *Err) String() string {
// 	return e.Error()
// }

var errNil = errors.New("nil")

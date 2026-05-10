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

func WrapError(e *Err, err error) error {
	e.err = err
	return e
}

type ErrType string

const (
	ErrTypeRequestValidation ErrType = "ERR_TYPE_REQUEST_VALIDATION"
	ErrTypeStorage           ErrType = "ERR_TYPE_STORAGE"
	ErrTypeBusinessLogic     ErrType = "ERR_TYPE_BUSINESS_LOGIC"
	ErrInternal              ErrType = "ERR_TYPE_INTERNAL"
)

func (e *Err) Error() string {
	if e == nil {
		return errNil.Error()
	}

	if e.err == nil {
		e.err = errNil
	}

	return fmt.Sprintf("%s: %s", e.String(), e.err.Error())
}

func (e *Err) Unwrap() error {
	return e.err
}

func (e *Err) String() string {
	if e == nil {
		return errNil.Error()
	}

	return fmt.Sprintf("error in %s.%s of type [%s]", e.SvcName, e.Method, e.Type)
}

var errNil = errors.New("nil")

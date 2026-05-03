package service

import "fmt"

type Err struct {
	SvcName string
	Method  string
	Type    ErrType
	Err     error
}

type ErrType string

const (
	ErrTypeRequestValidation ErrType = "ERR_TYPE_REQUEST_VALIDATION"
	ErrTypeStorage           ErrType = "ERR_TYPE_STORAGE"
	ErrTypeBusinessLogic     ErrType = "ERR_TYPE_BUSINESS_LOGIC"
)

func (e *Err) Error() string {
	return fmt.Sprintf("error in %s.%s of type [%s]:%s", e.SvcName, e.Method, e.Type, e.Err.Error())
}

func (e *Err) Unwrap() error {
	return e.Err
}

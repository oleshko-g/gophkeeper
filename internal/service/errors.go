package service

type Err struct {
	SvcName string
	Method  string
	Err     error
}

func (e *Err) Error() string {
	return e.Error()
}

func (e *Err) Unwrap() error {
	return e.Err
}

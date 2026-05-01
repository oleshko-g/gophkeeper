package service

type Err struct {
	SvcName string
	Method  string
	Err     error
}

func (e *Err) Error() string {
	return "error " + e.SvcName + "." + e.Method + ": " + e.Err.Error()
}

func (e *Err) Unwrap() error {
	return e.Err
}

package rpcutil

type UnexpectedError struct {
	Err error
}

func (e *UnexpectedError) Error() string {
	return "UnexpectedError: " + e.Err.Error()
}

func UnexpectedIfError(err error) error {
	if err == nil {
		return nil
	}
	return &UnexpectedError{err}
}

func IsUnexpectedIfError(err error) bool {
	if _, ok := err.(*UnexpectedError); ok {
		return true
	}
	return false
}

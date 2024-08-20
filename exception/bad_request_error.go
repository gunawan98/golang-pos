package exception

type BadRequestError struct {
	Error string
}

func NewBadRequestError(message string) BadRequestError {
	return BadRequestError{
		Error: message,
	}
}

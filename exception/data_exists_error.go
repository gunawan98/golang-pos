package exception

type DataAlreadyExistsError struct {
	Error string
}

func NewDataAlreadyExistsError(message string) DataAlreadyExistsError {
	return DataAlreadyExistsError{
		Error: message,
	}
}

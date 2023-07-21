package handlers

type GeneralError struct {
	Message string `json:"message"`
}

func (e GeneralError) Error() string {
	return e.Message
}

func handleError(err error) error {
	return GeneralError{Message: err.Error()}
}

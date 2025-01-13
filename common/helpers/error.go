package helpers

type ServiceError struct {
	Message string
	Code    int
}

func (s ServiceError) Error() string {
	return s.Message
}

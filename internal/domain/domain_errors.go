package domain

type DomainError struct {
	Cod     string
	Message string
}

func NewDomainError(
	cod string,
	message string,
) DomainError {

	return DomainError{
		Cod:     cod,
		Message: message,
	}
}

func (dError DomainError) GetCod() string {
	return dError.Cod
}

func (dError DomainError) Error() string {
	return dError.Message
}

var (
	ErrClientNotFound       = NewDomainError("client_not_found", "client not found")
	ErrClientWithoutBalance = NewDomainError("client_without_balance", "client without balance")
)

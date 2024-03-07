package interfaces

type Service interface {
	Execute(inputData interface{}) (output interface{}, err error)
}

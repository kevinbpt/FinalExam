package model

type DataResponse[T any] struct {
	Status  int
	Message string
	Data    T
}

type DataResponses[T any] struct {
	Status  int
	Message string
	Data    []T
}

type Message struct {
	Message string
}

type IError struct {
	Field   string
	Message string
	Value   string
}

type IValidate struct {
	CustomStruct interface{}
}

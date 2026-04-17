package storage

//go:generate moq -out depositor_mock.go . Depositor
type Depositor interface {
}

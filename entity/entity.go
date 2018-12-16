package entity

// A uuid that will always refer to an entity in the state
type ID string

type Entity interface {
	Height() int
	ID() ID
}

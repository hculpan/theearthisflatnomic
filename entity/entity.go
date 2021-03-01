package entity

// Entity represents the generic
// type of db entity
type Entity interface {
	Insert() error
	Save() error
}

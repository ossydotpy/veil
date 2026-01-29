package store

import "errors"

var (
	ErrNotFound = errors.New("secret not found")
)

type Store interface {
	Save(vault, name, value string) error
	Get(vault, name string) (string, error)
	Delete(vault, name string) error
	List(vault string) ([]string, error)
	ListVaults() ([]string, error)
	Nuke() error
	Close() error
}

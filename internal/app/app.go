package app

import (
	"github.com/ossydotpy/veil/internal/crypto"
	"github.com/ossydotpy/veil/internal/store"
)

type App struct {
	store  store.Store
	crypto *crypto.Engine
}

func New(s store.Store, c *crypto.Engine) *App {
	return &App{
		store:  s,
		crypto: c,
	}
}

func (a *App) Set(vault, name, value string) error {
	encrypted, err := a.crypto.Encrypt(value)
	if err != nil {
		return err
	}
	return a.store.Save(vault, name, encrypted)
}

func (a *App) Get(vault, name string) (string, error) {
	encrypted, err := a.store.Get(vault, name)
	if err != nil {
		return "", err
	}
	return a.crypto.Decrypt(encrypted)
}

func (a *App) Delete(vault, name string) error {
	return a.store.Delete(vault, name)
}

func (a *App) List(vault string) ([]string, error) {
	return a.store.List(vault)
}

func (a *App) ListVaults() ([]string, error) {
	return a.store.ListVaults()
}

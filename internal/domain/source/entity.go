package source

import (
	"crypto/rand"
	"encoding/hex"
)

type Source struct {
	ID     string
	Name   string
	Secret string
}

func NewSource(name string) *Source {
	return &Source{
		ID:     generateID(),
		Name:   name,
		Secret: generateSecret(),
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSecret() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Source) Verify(id string, secret string) bool {
	return s.ID == id && s.Secret == secret
}

func (s *Source) ResetSecret() {
	s.Secret = generateSecret()
}
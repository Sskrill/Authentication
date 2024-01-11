package hash

import (
	"crypto/sha1"
	"fmt"
)

type PSHasher interface {
	Hash(password string) (string, error)
}

type SHA1Hash struct {
	salt string
}

func NewPSHasher(salt string) *SHA1Hash {
	return &SHA1Hash{salt: salt}
}

func (s *SHA1Hash) Hash(password string) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(password)); err != nil {

		return "", nil
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt))), nil
}

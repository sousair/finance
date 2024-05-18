package cipher

import "golang.org/x/crypto/bcrypt"

type Cipher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type Bcrypt struct {
	hashCost int
}

var _ Cipher = (*Bcrypt)(nil)

func NewBcrypt(hashCost int) *Bcrypt {
	return &Bcrypt{hashCost: hashCost}
}

func (b Bcrypt) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.hashCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (b Bcrypt) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

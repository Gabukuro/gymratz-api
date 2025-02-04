package user

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type (
	Model struct {
		bun.BaseModel `bun:"users"`
		base.Model
		Name     string `bun:"name"`
		Email    string `bun:"email"`
		Password string `bun:"password"`
	}
)

func (m *Model) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	m.Password = string(bytes)

	return err
}

func (m *Model) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return err == nil
}

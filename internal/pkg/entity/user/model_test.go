package user_test

import (
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	t.Run("should hash the password and overwrite the struct field with the hashed value", func(t *testing.T) {
		t.Parallel()

		testPassword := "plain_password"
		model := &user.Model{
			Password: testPassword,
		}

		err := model.HashPassword()
		assert.Nil(t, err)
		assert.NotEqual(t, testPassword, model.Password)

		err = bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(testPassword))
		assert.Nil(t, err)
	})
}

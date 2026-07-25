package errorsmap

import (
	"errors"
	"strings"
	"testing"

	"github.com/josuebrunel/gopkg/assert"
)

func TestErrorsMap(t *testing.T) {

	const (
		userNotFound = "user-not-found"
		emailInvalid = "email-invalid"
	)
	var (
		em = New()
	)

	t.Run("Nil", func(t *testing.T) {
		assert.Eq(t, em.Nil(), true)
	})
	t.Run("ErrorOnEmptyMap", func(t *testing.T) {
		assert.Eq(t, em.Error(), "")
	})

	em["user"] = errors.New(userNotFound)
	t.Run("Get", func(t *testing.T) {
		assert.Eq(t, em.Get("toto"), "")
		assert.Eq(t, em.Get("user"), userNotFound)
		assert.Eq(t, em.Nil(), false)
	})
	t.Run("IfNil", func(t *testing.T) {
		assert.Eq(t, em.IfNil("user"), false)
		assert.Eq(t, em.IfNil("email"), true)
	})
	em["email"] = errors.New(emailInvalid)
	t.Run("Error", func(t *testing.T) {
		assert.Eq(t, strings.Contains(em.Error(), userNotFound), true)
		assert.Eq(t, strings.Contains(em.Error(), emailInvalid), true)
		assert.Eq(t, strings.Contains(em.Error(), "test"), false)
	})
}

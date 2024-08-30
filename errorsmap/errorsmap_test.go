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
		assert.AssertT(t, em.Nil(), true)
	})

	em["user"] = errors.New(userNotFound)
	t.Run("Get", func(t *testing.T) {
		assert.AssertT(t, em.Get("toto"), "")
		assert.AssertT(t, em.Get("user"), userNotFound)
		assert.AssertT(t, em.Nil(), false)
	})
	t.Run("IfNil", func(t *testing.T) {
		assert.AssertT(t, em.IfNil("user"), false)
		assert.AssertT(t, em.IfNil("email"), true)
	})
	em["email"] = errors.New(emailInvalid)
	t.Run("Error", func(t *testing.T) {
		assert.AssertT(t, strings.Contains(em.Error(), userNotFound), true)
		assert.AssertT(t, strings.Contains(em.Error(), emailInvalid), true)
		assert.AssertT(t, strings.Contains(em.Error(), "test"), false)
	})
}

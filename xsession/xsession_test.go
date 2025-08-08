package xsession

import (
	"testing"

	"github.com/josuebrunel/gopkg/assert"
)

func TestXUser(t *testing.T) {
	u := NewXUser("1", "token", "email", map[string]string{})
	u.SetData("plan", "free")
	assert.Eq(t, u.GetData("plan"), "free")
	u.DeleteData("plan")
	assert.Eq(t, u.GetData("plan"), "")
}

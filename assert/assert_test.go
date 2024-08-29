package assert

import (
	"testing"

	"golang.org/x/exp/constraints"
)

type (
	TestCase[T comparable, V comparable] struct {
		A, B T
		R    V
	}
	User struct {
		ID   int
		Name string
	}
)

func Add[T constraints.Ordered](a, b T) T {
	return a + b
}

var (
	uu = []User{
		{ID: 1, Name: "User1"},
		{ID: 2, Name: "User2"},
		{ID: 1, Name: "User1"},
	}
	tcInt   = []TestCase[int, int]{{1, 2, 3}, {10, 2, 12}}
	tcStr   = []TestCase[string, string]{{"hello", "world", "helloworld"}, {"pocket", "base", "pocketbase"}}
	tcUsers = []TestCase[User, bool]{
		{uu[0], uu[2], true},
		{uu[1], uu[2], false},
	}
)

func TestAdd(t *testing.T) {
	t.Run("intWithAssert", func(t *testing.T) {
		for _, tc := range tcInt {
			Assert(t, Add(tc.A, tc.B), tc.R)
		}
	})
	t.Run("intWithAssertT", func(t *testing.T) {
		for _, tc := range tcInt {
			Assert(t, Add(tc.A, tc.B), tc.R)
		}
	})
	t.Run("strWithAssert", func(t *testing.T) {
		for _, tc := range tcStr {
			Assert(t, Add(tc.A, tc.B), tc.R)
		}
	})
	t.Run("strWithAssertT", func(t *testing.T) {
		for _, tc := range tcStr {
			Assert(t, Add(tc.A, tc.B), tc.R)
		}
	})
	t.Run("userWithAssertT", func(t *testing.T) {
		for _, tc := range tcUsers {
			AssertT(t, tc.A == tc.B, tc.R)
		}
	})
	t.Run("userWithIn", func(t *testing.T) {
		AssertT(t, In(uu, uu[2]), true)
		AssertT(t, In(uu, User{7, "User7"}), false)
	})
}

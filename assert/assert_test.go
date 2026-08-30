package assert

import (
	"cmp"
	"testing"
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

func Add[T cmp.Ordered](a, b T) T {
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

func TestEq(t *testing.T) {
	t.Run("Ints", func(t *testing.T) {
		for _, tc := range tcInt {
			Eq(t, Add(tc.A, tc.B), tc.R)
		}
	})
	t.Run("Strs", func(t *testing.T) {
		for _, tc := range tcStr {
			Eq(t, Add(tc.A, tc.B), tc.R)
		}
	})
	t.Run("Users", func(t *testing.T) {
		for _, tc := range tcUsers {
			Eq(t, tc.A == tc.B, tc.R)
		}
	})
}

func TestNotEq(t *testing.T) {
	NotEq(t, 1, 2)
	NotEq(t, "a", "b")
}

func TestTrueFalse(t *testing.T) {
	True(t, 1 == 1)
	False(t, 1 == 2)
}

func TestNilNotNil(t *testing.T) {
	var (
		p *User
		m map[string]int
		s []int
	)
	Nil(t, nil)
	Nil(t, p)
	Nil(t, m)
	Nil(t, s)
	NotNil(t, &User{})
	NotNil(t, 1)
	NotNil(t, "x")
}

func TestContainsNotContains(t *testing.T) {
	Contains(t, uu, uu[2])
	NotContains(t, uu, User{7, "User7"})
}

// fakeT is a minimal assert.TestingT recorder used to verify that an
// assertion reports failure without tripping the real *testing.T running
// this test (a real t.Errorf would mark this test itself as failed).
type fakeT struct{ failed bool }

func (f *fakeT) Helper()                        {}
func (f *fakeT) Errorf(format string, a ...any) { f.failed = true }

// TestFailurePaths verifies every assertion returns false and reports a
// failure on mismatched input.
func TestFailurePaths(t *testing.T) {
	check := func(t *testing.T, name string, fn func(t TestingT) bool) {
		t.Helper()
		f := &fakeT{}
		got := fn(f)
		if got || !f.failed {
			t.Fatalf("%s: want (false, reported), got (%v, %v)", name, got, f.failed)
		}
	}

	check(t, "Eq", func(t TestingT) bool { return Eq(t, 1, 2) })
	check(t, "NotEq", func(t TestingT) bool { return NotEq(t, 1, 1) })
	check(t, "True", func(t TestingT) bool { return True(t, false) })
	check(t, "False", func(t TestingT) bool { return False(t, true) })
	check(t, "Nil", func(t TestingT) bool { return Nil(t, 1) })
	check(t, "NotNil", func(t TestingT) bool { return NotNil(t, nil) })
	check(t, "Contains", func(t TestingT) bool { return Contains(t, uu, User{7, "User7"}) })
	check(t, "NotContains", func(t TestingT) bool { return NotContains(t, uu, uu[2]) })
}

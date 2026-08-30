// Package assert provides minimal, dependency-free, non-fatal assertion
// helpers for use in tests. Each function returns true if the assertion
// passed; callers that need fail-fast behavior should check the return
// value and call t.FailNow() themselves.
package assert

import "reflect"

// TestingT is the subset of *testing.T (and *testing.B) that the
// assertions in this package depend on. It lets tests exercise these
// assertions against a fake recorder instead of a real *testing.T, whose
// Errorf would otherwise mark the enclosing test as failed.
type TestingT interface {
	Helper()
	Errorf(format string, args ...any)
}

// Eq asserts that got equals want.
func Eq[T comparable](t TestingT, got, want T) bool {
	t.Helper()
	if got != want {
		t.Errorf("assert.Eq: got %v, want %v", got, want)
		return false
	}
	return true
}

// NotEq asserts that got does not equal notWant.
func NotEq[T comparable](t TestingT, got, notWant T) bool {
	t.Helper()
	if got == notWant {
		t.Errorf("assert.NotEq: got %v, want anything but %v", got, notWant)
		return false
	}
	return true
}

// True asserts that got is true.
func True(t TestingT, got bool) bool {
	t.Helper()
	if !got {
		t.Errorf("assert.True: got false, want true")
		return false
	}
	return true
}

// False asserts that got is false.
func False(t TestingT, got bool) bool {
	t.Helper()
	if got {
		t.Errorf("assert.False: got true, want false")
		return false
	}
	return true
}

// Nil asserts that got is nil, treating typed nil pointers, maps, slices,
// chans, funcs, and interfaces as nil.
func Nil(t TestingT, got any) bool {
	t.Helper()
	if !isNil(got) {
		t.Errorf("assert.Nil: got %v, want nil", got)
		return false
	}
	return true
}

// NotNil asserts that got is not nil.
func NotNil(t TestingT, got any) bool {
	t.Helper()
	if isNil(got) {
		t.Errorf("assert.NotNil: got nil, want non-nil")
		return false
	}
	return true
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

// Contains asserts that list contains item.
func Contains[T comparable](t TestingT, list []T, item T) bool {
	t.Helper()
	if !contains(list, item) {
		t.Errorf("assert.Contains: %v not found in %v", item, list)
		return false
	}
	return true
}

// NotContains asserts that list does not contain item.
func NotContains[T comparable](t TestingT, list []T, item T) bool {
	t.Helper()
	if contains(list, item) {
		t.Errorf("assert.NotContains: %v found in %v, want absent", item, list)
		return false
	}
	return true
}

func contains[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

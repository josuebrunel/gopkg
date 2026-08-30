# assert

`assert` is a minimal, dependency-free assertion library for Go tests. Every assertion is non-fatal: it reports a failure via `t.Errorf` and returns `bool` (`true` on success), so a table-driven test keeps running and reports every failing case in one run instead of stopping at the first.

## Features

- **No dependencies**: pure stdlib (`reflect`), nothing else.
- **Generic**: `Eq`, `NotEq`, `Contains`, and `NotContains` work with any `comparable` type.
- **Non-fatal by default**: uses `t.Errorf`, not `t.Fatalf`. If you need fail-fast behavior, check the returned `bool` and call `t.FailNow()` yourself.
- **Works with `*testing.T` and `*testing.B`**: functions take the `assert.TestingT` interface, which both satisfy.
- **Correct nil handling**: `Nil`/`NotNil` use `reflect` to correctly treat typed-nil pointers, maps, slices, chans, and funcs as nil, not just `interface{} == nil`.

## Installation

```bash
go get github.com/josuebrunel/gopkg/assert
```

## Usage

```go
package main

import (
	"testing"

	"github.com/josuebrunel/gopkg/assert"
)

func TestExample(t *testing.T) {
	assert.Eq(t, 2+2, 4)
	assert.NotEq(t, "foo", "bar")
	assert.True(t, len("abc") == 3)
	assert.False(t, len("abc") == 4)
	assert.Nil(t, nil)
	assert.NotNil(t, &struct{}{})
	assert.Contains(t, []int{1, 2, 3}, 2)
	assert.NotContains(t, []int{1, 2, 3}, 4)
}
```

### Fail-fast

Every assertion returns `true` on success, `false` on failure, so you can bail out of a test early when continuing would be unsafe (e.g. before dereferencing a value you just checked for nil):

```go
func TestUser(t *testing.T) {
	u := loadUser()
	if !assert.NotNil(t, u) {
		t.FailNow()
	}
	assert.Eq(t, u.Name, "alice")
}
```

### Table-driven tests

Because assertions don't stop the test, every row in a table is checked and reported, not just the first failure:

```go
func TestAdd(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{1, 2, 3},
		{2, 2, 5}, // deliberately wrong, to show it doesn't block the next case
		{3, 3, 6},
	}
	for _, tc := range cases {
		assert.Eq(t, tc.a+tc.b, tc.want)
	}
}
```

### Testing with a fake `TestingT`

Because assertions depend only on the small `assert.TestingT` interface (`Helper()` + `Errorf(...)`), you can pass in your own recorder instead of a real `*testing.T` — useful when testing helpers that themselves call into `assert`:

```go
type fakeT struct{ failed bool }

func (f *fakeT) Helper()                        {}
func (f *fakeT) Errorf(format string, a ...any) { f.failed = true }

func TestMyHelperFails(t *testing.T) {
	f := &fakeT{}
	ok := assert.Eq(f, 1, 2)
	assert.False(t, ok)
	assert.True(t, f.failed)
}
```

## API

| Function                                      | Description                                  |
| ---------------------------------------------- | --------------------------------------------- |
| `Eq[T comparable](t, got, want T) bool`        | asserts `got == want`                          |
| `NotEq[T comparable](t, got, notWant T) bool`  | asserts `got != notWant`                       |
| `True(t, got bool) bool`                       | asserts `got` is `true`                        |
| `False(t, got bool) bool`                      | asserts `got` is `false`                       |
| `Nil(t, got any) bool`                         | asserts `got` is nil (typed-nil aware)         |
| `NotNil(t, got any) bool`                      | asserts `got` is not nil (typed-nil aware)     |
| `Contains[T comparable](t, list []T, item T) bool` | asserts `item` is in `list`               |
| `NotContains[T comparable](t, list []T, item T) bool` | asserts `item` is not in `list`        |

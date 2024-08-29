package assert

import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, x, y any) {
	if !reflect.DeepEqual(x, y) {
		t.Fatalf("[ASSERT-FAILED] - %v != %v", x, y)
	}
}

func AssertT[T comparable](t *testing.T, x, y T) {
	t.Helper()
	if x != y {
		t.Fatalf("[ASSERT-FAILED] - %v != %v", x, y)
	}
}

func In[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

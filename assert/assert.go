package assert

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Assert(t *testing.T, x, y any) {
	_, f, l, _ := runtime.Caller(1)
	f = filepath.Base(f)
	if !reflect.DeepEqual(x, y) {
		log.Fatalf("(%s:%d)[ASSERT-FAILED] - %v != %v", f, l, x, y)
	}
}

func AssertT[T comparable](t *testing.T, x, y T) {
	_, f, l, _ := runtime.Caller(1)
	f = filepath.Base(f)
	if x != y {
		log.Fatalf("(%s:%d)[ASSERT-FAILED] - %v != %v", f, l, x, y)
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

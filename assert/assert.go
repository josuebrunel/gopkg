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
		log.Fatalf("(%s:%d)[ASSERT-ERROR] - %v != %v", f, l, x, y)
	}
}

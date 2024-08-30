package errorsmap

import "errors"

type EMap map[string]error

func (e EMap) Nil() bool {
	for _, v := range e {
		if v != nil {
			return false
		}
	}
	return true
}

func (e EMap) IfNil(key string) bool {
	return e[key] == nil
}

func (e EMap) Get(key string) string {
	if v, ok := e[key]; ok {
		return v.Error()
	}
	return ""
}

func (e EMap) Error() string {
	errs := make([]error, 0)
	for _, v := range e {
		errs = append(errs, v)
	}
	return errors.Join(errs...).Error()
}

func New() EMap {
	return make(EMap)
}

package etcdutil

import "errors"



var ErrFieldNotFount = errors.New("field is not found")

var ErrGivenNotAPointer = errors.New("given is not a pointer")

var ErrGivenTypeNotMatch = errors.New("given type is not match")

var ErrGivenIdNotMatch = errors.New("given id is not match")

var ErrFieldNotSettable = errors.New("field is not settable")
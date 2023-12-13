package util

import (
	"bytes"
	"encoding/gob"
	"reflect"
)





func Deepcopy(src, dst interface{}) error {
	dst = DefaultStruct(src)
	gob.Register(src)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dst)
}


func DefaultStruct(src interface{}) interface{} {

	t := reflect.TypeOf(src)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	i := reflect.New(t)
	return i.Interface()
}



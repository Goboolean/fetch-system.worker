package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


type Model struct {
	Name string
	Age  int
}


var cases []struct {
	name string
	_default interface{}
	data interface{}
} = []struct{
	name string
	_default interface{}
	data interface{}
}{
	{
		name: "Model",
		_default: &Model{},
		data: &Model{
			Name: "goboolean",
			Age: 10,
		},
	},
}



func Test_Deepcopy(t *testing.T) {

	t.Skip("Deepcopy is not implemented yet")
	
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var dst interface{}
			err := Deepcopy(tt.data, dst)
			assert.NoError(t, err)
			assert.Equal(t, tt.data, dst)
		})
	}
}


func Test_DefaultStruct(t *testing.T) {
	
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			dst := DefaultStruct(tt.data)
			assert.Equal(t, tt._default, dst)
		})
	}
}
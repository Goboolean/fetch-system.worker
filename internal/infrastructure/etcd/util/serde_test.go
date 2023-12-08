package etcdutil_test

import (
	"reflect"
	"testing"

	etcdutil "github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd/util"
	"github.com/Goboolean/fetch-system.worker/internal/util"
	"github.com/stretchr/testify/assert"
)

func Contains[T any](list []T, target T) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, target) {
			return true
		}
	}
	return false
}

func Test_GroupByPrefix(t *testing.T) {

	type args struct {
		list   map[string]string
		prefix string
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := etcdutil.GroupByPrefix(tt.whole)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.group), len(got))
			for _, v := range tt.group {
				assert.True(t, Contains(got, v.str))
			}
		})
	}
}

func Test_Serialize(t *testing.T) {

	for _, tt := range cases {
		for _, ttt := range tt.group {
			t.Run(tt.name, func(t *testing.T) {
				str, err := etcdutil.Serialize(ttt.data)
				assert.NoError(t, err)
				assert.Equal(t, len(ttt.str), len(str))
				assert.Equal(t, ttt.str, str)
				assert.True(t, reflect.DeepEqual(ttt.str, str))
			})
		}
	}
}

func Test_Deserialize(t *testing.T) {

	for _, tt := range cases {
		for _, ttt := range tt.group {
			t.Run(tt.name, func(t *testing.T) {
				var input etcdutil.Model = util.DefaultStruct(ttt.model).(etcdutil.Model)
				err := etcdutil.Deserialize(ttt.str, input)
				assert.NoError(t, err)
				assert.Equal(t, ttt.data, input)
				assert.True(t, reflect.DeepEqual(ttt.data, input))
			})
		}
	}
}

func Test_SerializeDeserialize(t *testing.T) {

	for _, tt := range cases {
		for _, ttt := range tt.group {
			t.Run(tt.name, func(t *testing.T) {
				str, err := etcdutil.Serialize(ttt.data)
				assert.NoError(t, err)
				assert.Equal(t, len(ttt.str), len(str))
				assert.Equal(t, ttt.str, str)
				assert.True(t, reflect.DeepEqual(ttt.str, str))

				var input etcdutil.Model = util.DefaultStruct(ttt.model).(etcdutil.Model)
				err = etcdutil.Deserialize(str, input)
				assert.NoError(t, err)
				assert.Equal(t, ttt.data, input)
				assert.True(t, reflect.DeepEqual(ttt.data, input))
			})
		}
	}
}

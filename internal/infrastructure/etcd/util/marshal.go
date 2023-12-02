package etcdutil

import (
	"fmt"
	"reflect"
	"strings"

	"go.etcd.io/etcd/client/v3"
)

// GroupBy groups the list by the distinquisher that appears just after the prefix.
// It assums that list is sorted by key.
func GroupBy(list map[string]string) ([]map[string]string, error) {
	var m = make(map[string]map[string]string)
	var result []map[string]string

	var _type string
	for k := range list {
		p := strings.Split(k, "/")[1]
		if _type == "" {
			_type = p
		} else {
			if _type != p {
				return nil, ErrGivenTypeNotMatch
			}
		}
	}

	var prefix = Group(_type)

	for k, v := range list {
		p := strings.Split(strings.TrimPrefix(k, prefix), "/")[1]
		if _, ok := m[p]; !ok {
			m[p] = make(map[string]string)
		}
		m[p][k] = v
	}
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}



func Mmarshal(m Model) (map[string]string, error) {
	t := reflect.TypeOf(m)

	if t.Kind() != reflect.Ptr {
		return nil, ErrGivenNotAPointer
	}
	t = t.Elem()

	var result = make(map[string]string)
	var prefix string
	var id string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("etcd") == "" {
			continue
		}

		if field.Tag.Get("etcd") == "id" {
			id = reflect.ValueOf(m).Elem().FieldByName(field.Name).String()
			break
		}
	}

	prefix = fmt.Sprintf("/%s/%s/", m.Name(), id)
	result[strings.TrimSuffix(prefix, "/")] = ""

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("etcd")
		if key == "" || key == "id" {
			continue
		}

		var value = reflect.ValueOf(m).Elem().FieldByName(field.Name).String()
		result[prefix+key] = value
	}

	return result, nil
}

func MarshalList(list []Model) (map[string]string, error) {
	var result = make(map[string]string)

	for _, v := range list {
		m, err := Mmarshal(v)
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			result[k] = v
		}
	}

	return result, nil
}



func Unmarshal(str map[string]string, m Model) error {
	t := reflect.TypeOf(m)

	if t.Kind() != reflect.Ptr {
		return ErrGivenNotAPointer
	}
	t = t.Elem()


	var id string

	for k, v := range str {
		spl := strings.Split(k, "/")

		_type := spl[1]
		if _type != m.Name() {
			return ErrGivenTypeNotMatch
		}

		_id := spl[2]
		if id == "" {
			id = _id
		} else {
			if id != _id {
				return ErrGivenIdNotMatch
			}
		}

		if len(spl) == 3 {
			continue
		}

		key := spl[3]

		var name string

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Tag.Get("etcd") == key {
				name = field.Name
			}
		}

		f := reflect.ValueOf(m).Elem().FieldByName(name)
		if f.IsValid() == false || f.CanSet() == false {
			return ErrFieldNotSettable
		}

		f.SetString(v)
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("etcd") == "" {
			continue
		}

		if field.Tag.Get("etcd") == "id" {
			reflect.ValueOf(m).Elem().FieldByName(field.Name).SetString(id)
			break
		}
	}

	return nil
}


func Identifier(_type string, id string) string {
	return fmt.Sprintf("/%s/%s", _type, id)
}

func Group(_type string) string {
	return fmt.Sprintf("/%s", _type)
}

func Field(_type string, id string, field string) string {
	return fmt.Sprintf("/%s/%s/%s", _type, id, field)
}


func PayloadToMap(res *clientv3.GetResponse) map[string]string {
	var result = make(map[string]string)

	for _, v := range res.Kvs {
		result[string(v.Key)] = string(v.Value)
	}

	return result
}


// maybe unused: regarding to use formatting by t.(string)
func toStringFormat(m Model, f reflect.StructField) string {
	// TODO: fix float32 case problem
	switch f.Type.Kind() {
	case reflect.String:
		return reflect.ValueOf(m).FieldByName(f.Name).String()
	case reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64 | reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Float32:
		return fmt.Sprintf("%d", reflect.ValueOf(m).FieldByName(f.Name).Int())
	case reflect.Float64:
		return fmt.Sprintf("%f", reflect.ValueOf(m).FieldByName(f.Name).Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", reflect.ValueOf(m).FieldByName(f.Name).Bool())
	default:
		panic("not supported type")
	}
}
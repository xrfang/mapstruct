package mapstruct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Parse(mapVar, structVar interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	lookup := func(f reflect.StructField, m interface{}) (v interface{}, ok bool) {
		tag := f.Tag
		key := strings.ReplaceAll(tag.Get("yaml"), ",omitempty", "")
		if len(key) == 0 {
			key = strings.ReplaceAll(tag.Get("json"), ",omitempty", "")
		}
		if key == "" {
			key = strings.ToLower(f.Name)
		}
		switch m.(type) {
		case map[string]interface{}:
			v, ok = mapVar.(map[string]interface{})[key]
		case map[interface{}]interface{}:
			v, ok = mapVar.(map[interface{}]interface{})[key]
		default:
			panic(fmt.Errorf("invalid input data type (%T)", m))
		}
		return
	}
	get := func(f reflect.StructField, m interface{}) reflect.Value {
		ptr := f.Type.Kind() == reflect.Ptr
		v, ok := lookup(f, m)
		if ok {
			if ptr {
				return reflect.ValueOf(v).Convert(f.Type.Elem())
			}
			return reflect.ValueOf(v).Convert(f.Type)
		}
		return reflect.Zero(f.Type)
	}
	v := reflect.ValueOf(structVar)
	if v.Type().Kind() != reflect.Ptr {
		return errors.New("output data type must be a pointer")
	}
	v = v.Elem()
	if v.Type().Kind() != reflect.Struct {
		return errors.New("output data type must be a struct")
	}
	t := reflect.TypeOf(structVar).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		sf := t.Field(i)
		name := sf.Name
		kind := sf.Type.Kind()
		indi := false
		if kind == reflect.Ptr {
			kind = sf.Type.Elem().Kind()
			indi = true
		}
		switch kind {
		case reflect.Bool:
			v := get(sf, mapVar)
			if indi {
				b := v.Interface().(bool)
				v = reflect.ValueOf(&b)
			}
			f.Set(v)
		case reflect.Int:
			v := get(sf, mapVar)
			if indi {
				i := v.Interface().(int)
				v = reflect.ValueOf(&i)
			}
			f.Set(v)
		case reflect.Int8:
		case reflect.Int16:
		case reflect.Int32:
		case reflect.Int64:
		case reflect.Uint:
		case reflect.Uint8:
		case reflect.Uint16:
		case reflect.Uint32:
		case reflect.Uint64:
		case reflect.Uintptr:
		case reflect.Float32:
		case reflect.Float64:
			v := get(sf, mapVar)
			if indi {
				d := v.Interface().(float64)
				v = reflect.ValueOf(&d)
			}
			f.Set(v)
		case reflect.String:
			v := get(sf, mapVar)
			if indi {
				s := v.Interface().(string)
				v = reflect.ValueOf(&s)
			}
			f.Set(v)
		case reflect.Map:
		case reflect.Array:
		case reflect.Slice:
		case reflect.Struct:
			m, ok := lookup(sf, mapVar)
			if ok {
				var v reflect.Value
				if indi {
					v = reflect.New(sf.Type.Elem())
				} else {
					v = reflect.New(sf.Type)
				}
				assert(Parse(m, v.Interface()))
				if indi {
					f.Set(v)
				} else {
					f.Set(v.Elem())
				}
			}
		default:
			return fmt.Errorf("field `%s`: unsupported type (%v)", name, kind)
		}
	}
	return nil
}

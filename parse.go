package mapstruct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func getFieldTag(f reflect.StructField) string {
	tag := f.Tag
	key := strings.ReplaceAll(tag.Get("yaml"), ",omitempty", "")
	if len(key) == 0 {
		key = strings.ReplaceAll(tag.Get("json"), ",omitempty", "")
	}
	if key == "" {
		key = strings.ToLower(f.Name)
	}
	return key
}

func lookupMap(m interface{}, key string) (v interface{}, ok bool) {
	switch m.(type) {
	case map[string]interface{}:
		v, ok = m.(map[string]interface{})[key]
	case map[interface{}]interface{}:
		v, ok = m.(map[interface{}]interface{})[key]
	default:
		panic(fmt.Errorf("invalid input data type (%T)", m))
	}
	return
}

func getReflectValue(rt reflect.Type, v interface{}) (rv reflect.Value, err error) {
	kind := rt.Kind()
	indi := false
	if kind == reflect.Ptr {
		kind = rt.Elem().Kind()
		indi = true
	}
	if v == nil {
		rv = reflect.Zero(rt)
		return
	}
	rv = reflect.ValueOf(v)
	switch kind {
	case reflect.Bool:
		if indi {
			b := rv.Interface().(bool)
			rv = reflect.ValueOf(&b)
		}
	case reflect.Int:
		rv = rv.Convert(reflect.TypeOf(int(0)))
		if indi {
			d := rv.Interface().(int)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Int8:
		rv = rv.Convert(reflect.TypeOf(int8(0)))
		if indi {
			d := rv.Interface().(int8)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Int16:
		rv = rv.Convert(reflect.TypeOf(int16(0)))
		if indi {
			d := rv.Interface().(int16)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Int32:
		rv = rv.Convert(reflect.TypeOf(int32(0)))
		if indi {
			d := rv.Interface().(int32)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Int64:
		rv = rv.Convert(reflect.TypeOf(int64(0)))
		if indi {
			d := rv.Interface().(int64)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Uint:
		rv = rv.Convert(reflect.TypeOf(uint(0)))
		if indi {
			d := rv.Interface().(uint)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Uint8:
		rv = rv.Convert(reflect.TypeOf(uint8(0)))
		if indi {
			d := rv.Interface().(uint8)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Uint16:
		rv = rv.Convert(reflect.TypeOf(uint16(0)))
		if indi {
			d := rv.Interface().(uint16)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Uint32:
		rv = rv.Convert(reflect.TypeOf(uint32(0)))
		if indi {
			d := rv.Interface().(uint32)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Uint64:
		rv = rv.Convert(reflect.TypeOf(uint64(0)))
		if indi {
			d := rv.Interface().(uint64)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Uintptr:
		rv = rv.Convert(reflect.TypeOf(uintptr(0)))
		if indi {
			d := rv.Interface().(uintptr)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Float32:
		rv = rv.Convert(reflect.TypeOf(float32(0)))
		if indi {
			d := rv.Interface().(float32)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Float64:
		rv = rv.Convert(reflect.TypeOf(float64(0)))
		if indi {
			d := rv.Interface().(float64)
			rv = reflect.ValueOf(&d)
		}
	case reflect.String:
		if indi {
			d := rv.Interface().(string)
			rv = reflect.ValueOf(&d)
		}
	case reflect.Map:
		err = errors.New("not implemented")
	case reflect.Array:
		//err = errors.New("not implemented")
		t := reflect.ArrayOf(rv.Len(), rv.Type().Elem())
		fmt.Printf("type of array: %v (%T)\n", t, t)
		vr := reflect.New(t)
		// for i := 0; i < rv.Len(); i++ {
		// vr.Field(i).Set(rv.Field(i))
		// }
		rv = vr.Elem()
	case reflect.Slice:
		err = errors.New("not implemented")
	case reflect.Struct:
		if v != nil {
			if indi {
				rv = reflect.New(rt.Elem())
			} else {
				rv = reflect.New(rt)
			}
			assert(Parse(v, rv.Interface()))
			if !indi {
				rv = rv.Elem()
			}
		}
	default:
		err = fmt.Errorf("unsupported type (%v)", kind)
	}
	return
}

func Parse(mapVar, structVar interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
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
		mv, ok := lookupMap(mapVar, getFieldTag(sf))
		if !ok {
			continue
		}
		rv, err := getReflectValue(sf.Type, mv)
		if err != nil {
			fmt.Printf("getReflectValue(%s): %v\n", name, err)
			continue
		}
		f.Set(rv)
	}
	return nil
}

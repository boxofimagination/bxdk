package defaults

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	timeDurationType = reflect.TypeOf(time.Second)

	// ErrPassValue is returned if the caller pass a value instead of a pointer
	ErrPassValue = errors.New("must pass a pointer, not a value")

	// ErrNotStruct is returned if the caller pass a pointer of non struct
	ErrNotStruct = errors.New("must pass a pointer of struct")
)

// SetDefault set default value from struct tag: default
// for example:
// 	type A struct {
//		S string `default:"this is default"`
// 	}
func SetDefault(v interface{}) error {
	val := reflect.ValueOf(v)
	// prevent silent error, if a value is sent, the original value won't change
	if val.Kind() != reflect.Ptr {
		return ErrPassValue
	}
	indirect := reflect.Indirect(val)
	// prevent panic when call NumField()
	if indirect.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	numfield := indirect.NumField()
	for i := 0; i < numfield; i++ {
		fi := indirect.Field(i)
		if !fi.CanSet() {
			continue
		}

		// continue if it is not empty value
		if !reflect.DeepEqual(reflect.Zero(fi.Type()).Interface(), fi.Interface()) {
			continue
		}

		f := indirect.Type().Field(i)
		t := f.Tag.Get("default")
		// continue if default tag is not available
		if t == "" {
			continue
		}

		// for special types which have their own parser
		switch f.Type {
		case timeDurationType:
			n, err := time.ParseDuration(t)
			if err != nil {
				return err
			}
			fi.Set(reflect.ValueOf(n))
			continue
		}

		// for primitive types
		switch f.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			n, err := strconv.ParseInt(t, 10, 64)
			if err != nil {
				return err
			}
			fi.Set(reflect.ValueOf(n).Convert(f.Type))
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
			fi.Set(reflect.ValueOf(n).Convert(f.Type))
		case reflect.String:
			fi.Set(reflect.ValueOf(t).Convert(f.Type))
		}
	}
	return nil
}
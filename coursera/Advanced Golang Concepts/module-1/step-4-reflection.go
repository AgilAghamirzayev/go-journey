package main

import (
	"fmt"
	"reflect"
)

func SetField(target any, field string, value any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("SetField needs a non-nil pointer to a struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("SetField needs pointer to struct; got %s", rv.Kind())
	}

	fv := rv.FieldByName(field)
	if !fv.IsValid() {
		return fmt.Errorf("no such field %q", field)
	}
	if !fv.CanSet() {
		return fmt.Errorf("field %q is not settable (unexported?)", field)
	}

	in := reflect.ValueOf(value)
	if in.Type().AssignableTo(fv.Type()) {
		fv.Set(in)
		return nil
	}
	if in.Type().ConvertibleTo(fv.Type()) {
		fv.Set(in.Convert(fv.Type()))
		return nil
	}
	return fmt.Errorf("cannot assign %s to %s", in.Type(), fv.Type())
}

func main() {
	err := SetField("string", "Name", "Agil")
	if err != nil {
		fmt.Println(err)
	}
}

func CallIfExists(obj any, method string, args ...any) (results []any, called bool, err error) {
	v := reflect.ValueOf(obj)
	m := v.MethodByName(method)
	if !m.IsValid() {
		return nil, false, nil
	}
	in := make([]reflect.Value, len(args))
	for i, a := range args {
		in[i] = reflect.ValueOf(a)
	}
	out := m.Call(in)
	results = make([]any, len(out))
	for i, o := range out {
		results[i] = o.Interface()
	}
	return results, true, nil
}

package edit

import (
	"reflect"
)

var allEditTypes = make(map[string]any)

func FindType[T any](name string) T {

	// This is supposed to copy and create a new instance of the object
	val := reflect.New(reflect.TypeOf(allEditTypes[name]))
	return val.Elem().Interface().(T)
}

func RegisterType(name string, editType any) {
	allEditTypes[name] = editType
}

type Field struct {
	Name, Type string
}

func TypeFields(t any) []Field {

	reflectFields := reflect.VisibleFields(reflect.TypeOf(t))
	fields := []Field{}

	for _, rf := range reflectFields {

		fields = append(fields, Field{
			Name: rf.Name,
			Type: rf.Type.Name(),
		})
	}

	return fields
}

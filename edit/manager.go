package edit

import (
	"errors"
	"reflect"
)

type EditType interface {
	Label() string
}

type ReferenceType struct {
	Id   string `json:"_id"`
	Type string `json:"_type"`
}

var allEditTypes = make(map[string]EditType)

func AllTypes() []string {

	keys := make([]string, 0, len(allEditTypes))
	for k := range allEditTypes {
		keys = append(keys, k)
	}

	return keys
}

func GetType[T EditType](name string) (*T, error) {

	foundType := allEditTypes[name]

	if foundType == nil {
		return nil, errors.New("Unable to find type for " + name)
	}

	// This is supposed to copy and create a new instance of the object
	val := reflect.New(reflect.TypeOf(allEditTypes[name]))

	cast := val.Elem().Interface().(T)
	return &cast, nil
}

func RegisterType(name string, editType EditType) {

	// Validate edit type
	tType := reflect.TypeOf(editType)
	if tType == nil || tType.Kind() != reflect.Struct {
		panic("Cannot register non-struct type " + name)
	}

	allEditTypes[name] = editType
}

type Field struct {
	Name  string
	Type  reflect.Type
	Value reflect.Value
}

func TypeFields(t *EditType) ([]Field, error) {

	// Validate before reflecting
	tType := reflect.TypeOf(*t)
	if tType == nil || tType.Kind() != reflect.Struct {
		return nil, errors.New("Type must be of struct to read fields.")
	}

	reflectFields := reflect.VisibleFields(tType)
	fields := []Field{}

	// Based on stackoverflow post here:
	// https://stackoverflow.com/questions/63421976/panic-reflect-call-of-reflect-value-fieldbyname-on-interface-value

	// Get value of provided object
	tObjVal := reflect.Indirect(reflect.ValueOf(t).Elem().Elem())

	for _, rf := range reflectFields {

		fields = append(fields, Field{
			Name:  rf.Name,
			Type:  rf.Type,
			Value: tObjVal.FieldByName(rf.Name),
		})
	}

	return fields, nil
}

func UpdateField(t *EditType, field *Field, val any) {

	if !reflect.TypeOf(val).AssignableTo(field.Type) {
		return
	}

	newValue := reflect.ValueOf(val)

	tObjVal := reflect.ValueOf(t).Elem()
	tmp := reflect.New(tObjVal.Elem().Type()).Elem()

	tmp.Set(tObjVal.Elem())
	tmp.FieldByName(field.Name).Set(newValue)

	tObjVal.Set(tmp)
	field.Value = newValue
}

package edit

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

type EditType interface {
	Label() string
}

type ReferenceType struct {
	Id   string `json:"_id" mapstructure:"_id"`
	Type string `json:"_type" mapstructure:"_type"`
}

func (rt *ReferenceType) UnmarshalJSON(data []byte) error {

	var obj struct {
		Id   string `json:"_id" mapstructure:"_id"`
		Type string `json:"_type" mapstructure:"_type"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	*rt = ReferenceType{
		Id:   obj.Id,
		Type: obj.Type,
	}

	return nil
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
	Name    string
	Type    reflect.Type
	Value   reflect.Value
	RefType string
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

		isRefType := rf.Type.AssignableTo(reflect.TypeFor[ReferenceType]())
		var refType string
		if isRefType {
			refType = rf.Tag.Get("editRef")
		}

		fields = append(fields, Field{
			Name:    rf.Name,
			Type:    rf.Type,
			Value:   tObjVal.FieldByName(rf.Name),
			RefType: refType,
		})
	}

	return fields, nil
}

func ConvertFromStr(field *Field, val string) (reflect.Value, error) {

	kind := field.Type.Kind()

	switch kind {
	case reflect.String:
		return reflect.ValueOf(val), nil
	case reflect.Int:
		i, e := strconv.Atoi(val)
		if e != nil {
			return reflect.Value{}, nil
		}
		return reflect.ValueOf(i), nil
	case reflect.TypeFor[ReferenceType]().Kind():
		// For a reference type, we are assuming the value is a valid ID
		refT := ReferenceType{
			Id:   val,
			Type: field.RefType,
		}

		return reflect.ValueOf(refT), nil
	default:
		return reflect.Value{}, errors.New("Unsupported type!")
	}
}

func UpdateField(t *EditType, field *Field, val string) error {

	newValue, err := ConvertFromStr(field, val)

	if err != nil {
		return err
	}

	tObjVal := reflect.ValueOf(t).Elem()
	tmp := reflect.New(tObjVal.Elem().Type()).Elem()

	tmp.Set(tObjVal.Elem())
	tmp.FieldByName(field.Name).Set(newValue)

	tObjVal.Set(tmp)
	field.Value = newValue

	return nil
}

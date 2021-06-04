package graphql

import (
	"errors"
	"reflect"
)

type Value struct {
	// Check these before valType
	isVar  bool
	isNull bool
	isEnum bool

	// depending on this field the below is filled in
	// Supported: Int, Float64, String, Bool, Array, Map
	// Maybe we should rename Map to Struct everywhere
	valType reflect.Kind

	variable     string
	intValue     int
	floatValue   float64
	stringValue  string
	booleanValue bool
	enumValue    string
	listValue    []Value
	objectValue  Arguments

	// Set this value if the value might be used on multiple places and the graphql typename is known
	// When using this struct to set data and this field is defined you should check it
	qlTypeName *string
}

func (v *Value) SetToValueOfAndExpect(other Value, expect reflect.Kind) error {
	if other.valType != expect {
		return errors.New("Value expected to be of type " + expect.String())
	}
	v.SetToValueOf(other)
	return nil
}

func (v *Value) SetToValueOf(other Value) {
	v.valType = other.valType
	switch other.valType {
	case reflect.String:
		v.stringValue = other.stringValue
	case reflect.Int:
		v.intValue = other.intValue
	case reflect.Float64:
		v.floatValue = other.floatValue
	case reflect.Bool:
		v.booleanValue = other.booleanValue
	case reflect.Array:
		v.listValue = other.listValue
	case reflect.Map:
		v.objectValue = other.objectValue
	}
}

func makeStringValue(val string) Value {
	return Value{
		valType:     reflect.String,
		stringValue: val,
	}
}

func makeBooleanValue(val bool) Value {
	return Value{
		valType:      reflect.Bool,
		booleanValue: val,
	}
}

func makeIntValue(val int) Value {
	return Value{
		valType:  reflect.Int,
		intValue: val,
	}
}

func makeFloatValue(val float64) Value {
	return Value{
		valType:    reflect.Float64,
		floatValue: val,
	}
}

func makeEnumValue(val string) Value {
	return Value{
		isEnum:    true,
		enumValue: val,
	}
}

func makeNullValue() Value {
	return Value{
		isNull: true,
	}
}

func makeArrayValue(list []Value) Value {
	if list == nil {
		list = []Value{}
	}
	return Value{
		valType:   reflect.Array,
		listValue: list,
	}
}

func makeStructValue(keyValues Arguments) Value {
	if keyValues == nil {
		keyValues = Arguments{}
	}
	return Value{
		valType:     reflect.Map,
		objectValue: keyValues,
	}
}

func makeVariableValue(varName string) Value {
	return Value{
		variable: varName,
		isVar:    true,
	}
}
package operation

import (
	"strings"
)

type Type struct {
	Name         string
	DefaultValue string
	ReflectType  string
}

func PrimitiveTypes() []Type {
	return []Type{
		{"string", "\"\"", "reflect.String"},
		{"bool", "false", "reflect.Bool"},
		{"int", "0", "reflect.Int"},
		{"int8", "int8(0)", "reflect.Int8"},
		{"int16", "int16(0)", "reflect.Int16"},
		{"int32", "int32(0)", "reflect.Int32"},
		{"int64", "int64(0)", "reflect.Int64"},
		{"uint", "uint(0)", "reflect.Uint"},
		{"uint8", "uint8(0)", "reflect.Uint8"},
		{"uint16", "uint16(0)", "reflect.Uint16"},
		{"uint32", "uint32(0)", "reflect.Uint32"},
		{"uint64", "uint64(0)", "reflect.Uint64"},
		{"float32", "float32(0)", "reflect.Float32"},
		{"float64", "float64(0)", "reflect.Float64"},
	}
}

type TypeInput struct {
	Type        string
	TypeTitle   string
	ReflectType string
	Output      []*TypeOutput
}
type TypeOutput struct {
	ParentType  string
	Type        string
	TypeTitle   string
	ReflectType string
	Value       string
}

var data =  func() []*TypeInput {
	data := make([]*TypeInput, len(PrimitiveTypes()))
	for index, in := range PrimitiveTypes() {
		input := &TypeInput{
			Type:        in.Name,
			TypeTitle:   strings.Title(in.Name),
			ReflectType: in.ReflectType,
			Output:      make([]*TypeOutput, len(PrimitiveTypes())),
		}
		for indexOut, out := range PrimitiveTypes() {
			output := &TypeOutput{
				ParentType:  in.Name,
				Type:        out.Name,
				TypeTitle:   strings.Title(out.Name),
				ReflectType: out.ReflectType,
				Value:       out.DefaultValue,
			}
			input.Output[indexOut] = output
		}
		data[index] = input
	}
	return data
}()

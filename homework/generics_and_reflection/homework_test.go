package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(v any) string {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ""
	}

	result := make([]string, 0, val.NumField())
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("properties")
		if tag == "" || tag == "-" {
			continue
		}

		tagParts := strings.Split(tag, ",")
		tag = tagParts[0]
		omitEmpty := false
		for _, opt := range tagParts[1:] {
			if opt == "omitempty" {
				omitEmpty = true
				break
			}
		}
		fv := val.Field(i)
		if omitEmpty && isZeroValue(fv) {
			continue
		}
		var strVal string

		switch fv.Kind() {
		case reflect.String:
			strVal = fv.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strVal = strconv.FormatInt(fv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			strVal = strconv.FormatUint(fv.Uint(), 10)
		case reflect.Bool:
			strVal = strconv.FormatBool(fv.Bool())
		case reflect.Float32, reflect.Float64:
			strVal = strconv.FormatFloat(fv.Float(), 'f', -1, 64)
		default:
			continue
		}

		result = append(result, tag+"="+strVal)
	}

	return strings.Join(result, "\n")

}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	}
	return false
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

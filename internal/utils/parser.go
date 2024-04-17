package utils

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DefaultParser struct {}

func (p *DefaultParser) StructToBson(v interface{}) bson.M {
    doc := bson.M{}
	val := reflect.ValueOf(v)

	// if v is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// ensure the field is exportable
		if field.CanInterface() {
			jsonTags := strings.Split(typ.Field(i).Tag.Get("json"), ",")
			jsonTag := jsonTags[0]
			if field.Kind() == reflect.Ptr && !field.IsNil() {
				processPtrField(field, jsonTag, &doc, p)
			} else if field.Kind() == reflect.Struct {
				// recursive call for non-pointer nested structs
				doc[jsonTag] = p.StructToBson(field.Interface())
			} else if field.Kind() != reflect.Ptr {
				doc[jsonTag] = field.Interface()
			}
		}
	}

	return doc
}

func processPtrField(field reflect.Value, jsonTag string, doc *bson.M, p *DefaultParser) {
	if field.Elem().Type() == reflect.TypeOf(primitive.ObjectID{}) {
		if !field.Elem().IsZero() {
			(*doc)[jsonTag] = field.Elem().Interface()
		}
		return
	} 

	if field.Elem().Kind() == reflect.Slice || field.Elem().Kind() == reflect.Array {
		(*doc)[jsonTag] = processSliceOrArray(field.Elem(), p)
	} else if field.Elem().Kind() == reflect.Struct {
		(*doc)[jsonTag] = p.StructToBson(field.Elem().Interface())
	} else {
		(*doc)[jsonTag] = field.Elem().Interface()
	}
}

// processSliceOrArray handles both slices and arrays, extracting their elements and processing them as needed.
func processSliceOrArray(val reflect.Value, p *DefaultParser) []interface{} {
    var result []interface{}
    for j := 0; j < val.Len(); j++ {
        elem := val.Index(j)
        if elem.Kind() == reflect.Ptr && !elem.IsNil() && elem.Elem().Kind() == reflect.Struct {
            result = append(result, p.StructToBson(elem.Elem().Interface()))
        } else if elem.Kind() == reflect.Struct {
            result = append(result, p.StructToBson(elem.Interface()))
        } else {
            result = append(result, elem.Interface())
        }
    }

    return result
}
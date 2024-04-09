package utils

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type DefaultParser struct {}

func (p *DefaultParser) StructToBson(v interface{}) bson.M {
    doc := bson.M{}
    val := reflect.ValueOf(v)
    
    // Check if v is a pointer and get the value it points to
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        if field.CanInterface() {
            jsonTags := strings.Split(typ.Field(i).Tag.Get("json"), ",")
            jsonTag := jsonTags[0] // Correctly handling omitempty

            // Handle case where field is a pointer to a slice or array
            if field.Kind() == reflect.Ptr && !field.IsNil() && (field.Elem().Kind() == reflect.Slice || field.Elem().Kind() == reflect.Array) {
                sliceData := processSliceOrArray(field.Elem(), p)
                doc[jsonTag] = sliceData
            } else if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
                sliceData := processSliceOrArray(field, p)
                doc[jsonTag] = sliceData
            } else if field.Kind() == reflect.Ptr && !field.IsNil() {
                // Check for a nested struct and process recursively
                if field.Elem().Kind() == reflect.Struct {
                    doc[jsonTag] = p.StructToBson(field.Elem().Interface())
                } else {
                    doc[jsonTag] = field.Elem().Interface()
                }
            } else if field.Kind() == reflect.Struct {
                doc[jsonTag] = p.StructToBson(field.Interface())
            } else {
                doc[jsonTag] = field.Interface()
            }
        }
    }

    return doc
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
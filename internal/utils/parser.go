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
				if field.Elem().Kind() == reflect.Struct {
					// recursive call for nested structs
					doc[jsonTag] = p.StructToBson(field.Elem().Interface())
				} else {
					doc[jsonTag] = field.Elem().Interface()
				}
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
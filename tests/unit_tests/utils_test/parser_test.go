package u

import (
	"reflect"
	"testing"

	"github.com/GhostDrew11/vigor-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func TestStructToBson(t *testing.T) {
	parser := utils.DefaultParser{}

	tests := []struct {
		name     string
		input    interface{}
		expected bson.M
	}{
		{
			name: "simple struct",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "John",
				Age:  25,
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
			},
		},
		{
			name : "simple struct with more than one json tag",
			input: struct {
				Name string `json:"name,omitempty"`
				Age  int    `json:"age,omitempty"`
			}{
				Name: "John",
				Age:  25,
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
			},
		},
		{
			name: "simple struct with pointer",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Email *string `json:"email"`
			}{
				Name: "John",
				Age:  25,
				Email: nil,
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
				// email is nil, so it should not be included in the output
			},
		},
		{
			name : "simple struct with non-empty pointer",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Email *string `json:"email"`
			}{
				Name: "John",
				Age:  25,
				Email: StringPtr("johndoe@example.com"),
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
				"email": "johndoe@example.com",
			},
		},
		{
			name: "nested struct",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Email *string `json:"email"`
				Address struct {
					Street string `json:"street"`
					City   string `json:"city"`
				} `json:"address"`
			}{
				Name: "John",
				Age:  25,
				Email: nil,
				Address:  struct {
					Street string `json:"street"`
					City   string `json:"city"`
				}{
					Street: "123 Main St",
					City:   "Springfield",
				},
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
				"address": bson.M{
					"street": "123 Main St",
					"city":   "Springfield",
				},
			},
		},
		{
			name: "nested struct with pointer",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Email *string `json:"email"`
				Address *struct {
					Street string `json:"street"`
					City   string `json:"city"`
				} `json:"address"`
			}{
				Name: "John",
				Age:  25,
				Email: nil,
				Address:  nil,
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
				// address is nil, so it should not be included in the output
			},
		},
		{
			name: "nested struct with non-empty pointer",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Email *string `json:"email"`
				Address *struct {
					Street string `json:"street"`
					City   string `json:"city"`
				} `json:"address"`
			}{
				Name: "John",
				Age:  25,
				Email: nil,
				Address: &struct {
					Street string `json:"street"`
					City   string `json:"city"`
				}{
					Street: "123 Main St",
					City:   "Springfield",
				},
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
				"address": bson.M{
					"street": "123 Main St",
					"city":   "Springfield",
				},
			},
		},
		{
			name: "nested struct with non-empty pointer and nested pointer",
			input: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				Email *string `json:"email"`
				Address *struct {
					Street string `json:"street"`
					City   string `json:"city"`
					Zip    *int   `json:"zip"`
				} `json:"address"`
			}{
				Name: "John",
				Age:  25,
				Email: nil,
				Address: &struct {
					Street string `json:"street"`
					City   string `json:"city"`
					Zip    *int   `json:"zip"`
				}{
					Street: "123 Main St",
					City:   "Springfield",
					Zip:    IntPtr(12345),
				},
			},
			expected: bson.M{
				"name": "John",
				"age":  25,
				"address": bson.M{
					"street": "123 Main St",
					"city":   "Springfield",
					"zip":    12345,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parser.StructToBson(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}


}
package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
)

// RouteHeader represents additional headers for a route
type RouteHeader struct {
	Key   string
	Value string
}

// Route represents the structure of a route in the API
type Route struct {
	Name        string
	Method      string
	Path        string
	Description string
	Headers     []RouteHeader
}

// PostmanItem represents a single item in a Postman collection
type PostmanItem struct {
	Name     string        `json:"name"`
	Request  Request       `json:"request"`
	Response []interface{} `json:"response"`
}

// Request represents the structure of a request in a Postman collection
type Request struct {
	Method string        `json:"method"`
	Header []interface{} `json:"header"`
	Body   Body          `json:"body"`
	URL    URL           `json:"url"`
}

// Body represents the body of the request in a Postman collection
type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

// URL represents the URL structure in Postman
type URL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

// PostmanCollection represents the top-level structure of a Postman collection
type PostmanCollection struct {
	Info Info          `json:"info"`
	Item []PostmanItem `json:"item"`
}

// Info represents metadata about the Postman collection
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}

func main() {
	// Define the routes
	routes := []Route{
		{
			Name:        "Register Admin",
			Method:      "POST",
			Path:        "/api/v1/auth/admin/register",
			Description: "Register a new admin user",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Login Admin",
			Method:      "POST",
			Path:        "/api/v1/auth/admin/login",
			Description: "Login as an admin user",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Register User",
			Method:      "POST",
			Path:        "/api/v1/auth/user/register",
			Description: "Register a new user",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Login User",
			Method:      "POST",
			Path:        "/api/v1/auth/user/login",
			Description: "Login as a user",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Refresh Token",
			Method:      "POST",
			Path:        "/api/v1/auth/refresh",
			Description: "Refresh the access token",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Create Exercise",
			Method:      "POST",
			Path:        "/api/v1/admin/exercises",
			Description: "Create a new exercise",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Get Exercises",
			Method:      "GET",
			Path:        "/api/v1/admin/exercises",
			Description: "Get all exercises",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Get Exercise by ID",
			Method:      "GET",
			Path:        "/api/v1/admin/exercises/:id",
			Description: "Get an exercise by its ID",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Search Exercises by Name",
			Method:      "GET",
			Path:        "/api/v1/admin/exercises/search",
			Description: "Search for exercises by name",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Update Exercise",
			Method:      "PUT",
			Path:        "/api/v1/admin/exercises/:id",
			Description: "Update an existing exercise",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
		{
			Name:        "Delete Exercise",
			Method:      "DELETE",
			Path:        "/api/v1/admin/exercises/:id",
			Description: "Delete an existing exercise",
			Headers: []RouteHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
	}

	// Attemp to read and update an existing collection; otherwise generate a new one
	postmanCollection, err := readAndUpdateCollection("postman_collection.json", routes)
	if err != nil {
		log.Fatalf("Error reading and updating Postman collection: %v", err)
	}

	// Marshal the Postman collection to JSON
	collectionJSON, err := json.MarshalIndent(postmanCollection, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling Postman collection: %v", err)
	}

	// Write JSON to file
	err = os.WriteFile("postman_collection.json", collectionJSON, 0644)
	if err != nil {
		log.Fatalf("Error writing Postman collection to file: %v", err)
	}

	log.Println("Postman collection generated successfully")
}

func readAndUpdateCollection(fileName string, newRoutes []Route) (PostmanCollection, error) {
	var collection PostmanCollection

	if _, err:= os.Stat(fileName); err == nil {
		file, readErr := os.ReadFile(fileName)
		if readErr != nil {
			return collection, readErr
		}

		if unmarshalErr := json.Unmarshal(file, &collection); unmarshalErr != nil {
			return collection, unmarshalErr
		}
	} else {
		collection = PostmanCollection{
			Info: Info{
				Name:        "Vigor API",
				Description: "Generated Postman collection for the Vigor API",
				Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
			},
			Item: []PostmanItem{},
		}
	}

	collection = updateCollectionWithRoutes(collection, newRoutes)
	return collection, nil
}

func updateCollectionWithRoutes(collection PostmanCollection, routes []Route) PostmanCollection {
	for _, newRoute := range routes {
		updated := false
		for i, item := range collection.Item {
			if item.Name == newRoute.Name {
				newItem := createPostmanItemFromRoute(newRoute)
				if !reflect.DeepEqual(item, newItem) {
					collection.Item[i] = newItem
					updated = true
				}
				break
			}
		}
		if !updated {
			collection.Item = append(collection.Item, createPostmanItemFromRoute(newRoute))
		}
	}

	return collection
}

func createPostmanItemFromRoute(route Route) PostmanItem {
	headers := make([]interface{}, 0, len(route.Headers))
	for _, header := range route.Headers {
		headers = append(headers, map[string]string{
			"key":   header.Key,
			"value": header.Value,
		})
	}

	routePath := route.Path
	routePath = strings.Replace(routePath, ":id", "{{exerciseId}}", -1)

	item := PostmanItem{
		Name: route.Name,
		Request: Request{
			Method: route.Method,
			Header: headers,
			Body: Body{
				Mode: "raw",
				Raw:  "{}",
			},
			URL: URL{
				Raw:  "{{baseUrl}}" + routePath,
				Host: []string{"{{baseUrl}}"},
				Path: []string{routePath},
			},
		},
		Response: []interface{}{},
	}

	return item
}



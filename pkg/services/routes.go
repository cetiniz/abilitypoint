package service

import (
	"net/http"
)

// Route defines how to form a new route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains array of server endpoints
type Routes []Route

var routes = Routes{
	Route{
		"GetAccount", // Name
		"GET",        // HTTP method
		"/api",       // Route pattern
		fetchGraph,
	},
}

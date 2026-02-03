package utils

import (
	"fmt"

	"github.com/gorilla/mux"
)

func PrintRoutes(r *mux.Router) {
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		fmt.Printf("ROUTE %-30s METHODS %v\n", path, methods)
		return nil
	})
}

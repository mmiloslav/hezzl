package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name               string
	Method             string
	Pattern            string
	HandlerFunc        http.HandlerFunc
	MiddlewareAuthFunc func(http.Handler) http.Handler
}

type Routes []Route

var routes = Routes{
	Route{Name: "Ping", Method: "GET", Pattern: "/api/ping", HandlerFunc: pingHandler, MiddlewareAuthFunc: emptyMiddleWare},
	Route{Name: "GoodCreate", Method: "POST", Pattern: "/api/good/create", HandlerFunc: goodCreate, MiddlewareAuthFunc: emptyMiddleWare},
	Route{Name: "GoodUpdate", Method: "PATCH", Pattern: "/api/good/update", HandlerFunc: goodUpdate, MiddlewareAuthFunc: emptyMiddleWare},
	Route{Name: "GoodDelete", Method: "DELETE", Pattern: "/api/good/delete", HandlerFunc: goodDelete, MiddlewareAuthFunc: emptyMiddleWare},
	Route{Name: "GoodsList", Method: "GET", Pattern: "/api/goods/list", HandlerFunc: goodsList, MiddlewareAuthFunc: emptyMiddleWare},
	Route{Name: "GoodReprioritize", Method: "PATCH", Pattern: "/api/good/reprioritize", HandlerFunc: goodReprioritize, MiddlewareAuthFunc: emptyMiddleWare},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.MiddlewareAuthFunc(route.HandlerFunc))
	}
	return router
}

func emptyMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(r.Context())
		handler.ServeHTTP(w, r)
	})
}

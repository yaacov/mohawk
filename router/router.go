package router

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Route struct {
	Method  string
	paths   []string
	handler func(http.ResponseWriter, *http.Request)
}

type Router struct {
	Prefix           string
	HandleBadRequest func(http.ResponseWriter, *http.Request)
	Routes           []Route
}

func (router *Router) Add(method string, path string, handler func(http.ResponseWriter, *http.Request)) {
	router.Routes = append(router.Routes, Route{method, strings.Split(path, "/"), handler})
}

func (_ Router) Find(route Route, paths []string, w http.ResponseWriter, r *http.Request) bool {
	// check method
	if r.Method != route.Method || len(paths) != len(route.paths) {
		return false
	}

	// check path
	for i, p := range route.paths {
		if p[0] != ':' && paths[i] != p {
			return false
		}
	}

	// get arguments
	for i, p := range route.paths {
		if p[0] == ':' {
			e, _ := url.QueryUnescape(paths[i])
			r.Form.Add(p[1:], e)
		}
	}

	// run handler
	return true
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %4s %s", r.RemoteAddr, r.Method, r.URL)

	// check prefix
	path := r.URL.EscapedPath()
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if !strings.HasPrefix(path, router.Prefix) {
		router.HandleBadRequest(w, r)
		return
	}

	// parse form
	r.ParseForm()

	// do route
	paths := strings.Split(path[len(router.Prefix):], "/")
	for _, route := range router.Routes {
		found := router.Find(route, paths, w, r)
		if found {
			route.handler(w, r)
			return
		}
	}

	// handle page not found
	router.HandleBadRequest(w, r)
}
package config

import (
	"net/http"
)

var global = new(Config)

func ReadAll() error {
	return global.ReadAll()
}

func Flags() map[string]any {
	return global.Flags()
}

func Env() map[string]any {
	return global.Env()
}

func Vars() []Var {
	return global.all
}

func HTTPHandler() http.Handler {
	return global.HandleHTTP()
}

package api

import (
	"github.com/gorilla/mux"
	"github.com/luisnquin/nao/src/data"
)

type StandardResponse struct {
	Version string            `json:"apiVersion"`
	Context string            `json:"context"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params,omitempty"`
	Error   any               `json:"error,omitempty"`
	Data    any               `json:"data"`
}

type Api struct {
	router *mux.Router
	box    *data.Box
}

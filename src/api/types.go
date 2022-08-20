package api

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/nao/src/store"
)

type StandardResponse struct {
	Version string `json:"apiVersion"`
	Context string `json:"context"`
	Method  string `json:"method"`
	Params  params `json:"params,omitempty"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"store"`
}

type Server struct {
	router       *echo.Echo
	box          *store.Box
	itWasMyFault chan bool
	port         string
	quiet        bool
}

type contentDTO struct {
	Content string `json:"content"`
}

type params map[string]string

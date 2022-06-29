package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/andreaskoch/go-fswatch"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/data"
)

func New() *Server {
	s := &Server{
		router: echo.New(),
		box:    data.New(),
	}

	s.router.Use(middleware.CORS())

	return s
}

func (a *Server) Start(port string) error {
	go a.listenAndRefreshData()
	a.mountHandlers()

	return http.ListenAndServe(port, a.router)
}

func (a *Server) listenAndRefreshData() {
	w := fswatch.NewFileWatcher(config.App.Paths.DataFile, 3)
	w.Start()

	for w.IsRunning() {
		select {
		case <-w.Modified():
			a.box.ModifyBox(data.JustLoadBox())
			color.New(color.FgHiBlue).Fprintln(os.Stdout, "Data refreshed")

		case <-w.Moved():
			color.New(color.FgHiRed).Fprintln(os.Stderr, "Error: Unable to find data file, apparently moved")
			os.Exit(1)
		}
	}
}

func (a *Server) mountHandlers() {
	sets := a.router.Group("/sets")
	sets.GET("", a.GetSetsHandler())
	sets.POST("", a.NewSetHandler())
	sets.GET("/:id", a.GetSetHandler())
	sets.PUT("/:id", a.ModifySetHandler())
	sets.DELETE("/:id", a.DeleteSetHandler())
	sets.PATCH("/:id", a.ModifySetContentHandler())
}

func (a *Server) JSONResponse(w http.ResponseWriter, statusCode int, v any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Del("Content-Type")
	}
}

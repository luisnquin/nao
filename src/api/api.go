package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/andreaskoch/go-fswatch"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/data"
)

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
		box:    data.New(),
	}
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
	a.router.HandleFunc("/sets", a.GetSetsHandler()).Methods(http.MethodGet)
	a.router.HandleFunc("/sets", a.NewSetHandler()).Methods(http.MethodPost)
	a.router.HandleFunc("/sets/{id}", a.GetSetHandler()).Methods(http.MethodGet)
	a.router.HandleFunc("/sets/{id}", a.ModifySetHandler()).Methods(http.MethodPut)
	a.router.HandleFunc("/sets/{id}", a.DeleteSetHandler()).Methods(http.MethodDelete)
	a.router.HandleFunc("/sets/{id}", a.ModifySetContentHandler()).Methods(http.MethodPatch)
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

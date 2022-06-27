package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
)

func (a *Server) GetSetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.JSONResponse(w, http.StatusOK, StandardResponse{
			Version: constants.Version,
			Method:  r.Method,
			Context: "sets",
			Data:    a.box.ListSets(),
		})
	}
}

func (a *Server) GetSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		set, err := a.box.GetSet(params["id"])
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		a.JSONResponse(w, http.StatusOK, StandardResponse{
			Version: constants.Version,
			Method:  r.Method,
			Context: "sets",
			Params:  params,
			Data:    set,
		})
	}
}

func (a *Server) NewSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var set data.Set

		err := json.NewDecoder(r.Body).Decode(&set)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		_ = r.Body.Close()

		k, err := a.box.NewFromSet(set)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		a.JSONResponse(w, http.StatusCreated, StandardResponse{
			Version: constants.Version,
			Method:  r.Method,
			Context: "sets",
			Data: map[string]string{
				"key": k,
			},
		})
	}
}

func (a *Server) ModifySetContentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			params  map[string]string = mux.Vars(r)
			request contentDTO
		)

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		err = a.box.ModifySetContent(params["id"], request.Content)
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *Server) ModifySetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			params  map[string]string = mux.Vars(r)
			request data.Set
		)

		err := a.box.OverwriteSet(params["id"], request)
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *Server) DeleteSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		err := a.box.DeleteSet(params["id"])
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

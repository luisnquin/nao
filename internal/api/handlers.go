package api

/*
import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/store"
)

func (a *Server) GetNotesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, StandardResponse{
			Version: config.Version,
			Method:  c.Request().Method,
			Context: "notes",
			Data:    a.box.List(),
		})
	}
}

func (a *Server) GetNoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		note, err := a.box.Get(c.Param("id"))
		if err != nil {
			if errors.Is(err, store.ErrNoteNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, StandardResponse{
			Version: config.Version,
			Method:  c.Request().Method,
			Context: "notes",
			Params: params{
				"id": c.Param("id"),
			},
			Data: note,
		})
	}
}

func (a *Server) NewNoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request store.Note

		err := c.Bind(&request)
		if err != nil {
			return echo.ErrBadRequest
		}

		k, err := a.box.NewFrom(request)
		if err != nil {
			return echo.ErrInternalServerError
		}

		a.itWasMyFault <- true

		return c.JSON(http.StatusCreated, StandardResponse{
			Version: config.Version,
			Method:  c.Request().Method,
			Context: "notes",
			Data: echo.Map{
				"key": k,
			},
		})
	}
}

func (a *Server) ModifyNoteContentHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request contentDTO

		err := c.Bind(&request)
		if err != nil {
			return echo.ErrBadRequest
		}

		err = a.box.ModifyContent(c.Param("id"), request.Content)
		if err != nil {
			if errors.Is(err, store.ErrNoteNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		a.itWasMyFault <- true

		return c.NoContent(http.StatusOK)
	}
}

func (a *Server) ModifyNoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request store.Note

		err := a.box.Replace(c.Param("id"), request)
		if err != nil {
			if errors.Is(err, store.ErrNoteNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		a.itWasMyFault <- true

		return c.NoContent(http.StatusOK)
	}
}

func (a *Server) DeleteNoteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := a.box.Delete(c.Param("id"))
		if err != nil {
			if errors.Is(err, store.ErrNoteNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		a.itWasMyFault <- true

		return c.NoContent(http.StatusOK)
	}
}
*/

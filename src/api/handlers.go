package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
)

func (a *Server) GetSetsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, StandardResponse{
			Version: constants.Version,
			Method:  c.Request().Method,
			Context: "sets",
			Data:    a.box.ListSets(),
		})
	}
}

func (a *Server) GetSetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		set, err := a.box.GetSet(c.Param("id"))
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		return c.JSON(http.StatusOK, StandardResponse{
			Version: constants.Version,
			Method:  c.Request().Method,
			Context: "sets",
			Params: params{
				"id": c.Param("id"),
			},
			Data: set,
		})
	}
}

func (a *Server) NewSetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request data.Set

		err := c.Bind(&request)
		if err != nil {
			return echo.ErrBadRequest
		}

		k, err := a.box.NewFromSet(request)
		if err != nil {
			return echo.ErrInternalServerError
		}

		a.itWasMe <- true

		return c.JSON(http.StatusCreated, StandardResponse{
			Version: constants.Version,
			Method:  c.Request().Method,
			Context: "sets",
			Data: echo.Map{
				"key": k,
			},
		})
	}
}

func (a *Server) ModifySetContentHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request contentDTO

		err := c.Bind(&request)
		if err != nil {
			return echo.ErrBadRequest
		}

		err = a.box.ModifySetContent(c.Param("id"), request.Content)
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		a.itWasMe <- true

		return c.NoContent(http.StatusOK)
	}
}

func (a *Server) ModifySetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request data.Set

		err := a.box.OverwriteSet(c.Param("id"), request)
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		a.itWasMe <- true

		return c.NoContent(http.StatusOK)
	}
}

func (a *Server) DeleteSetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := a.box.DeleteSet(c.Param("id"))
		if err != nil {
			if errors.Is(err, data.ErrSetNotFound) {
				return echo.ErrNotFound
			}

			return echo.ErrInternalServerError
		}

		a.itWasMe <- true

		return c.NoContent(http.StatusOK)
	}
}

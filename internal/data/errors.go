package data

import (
	"errors"
	"strings"
)

var ErrRunningOnHomelessShelter = errors.New("running on homeless shelter")

func isHomelessShelterError(err error) bool {
	return strings.Contains(err.Error(), "homeless-shelter")
}

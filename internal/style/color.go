package style

import (
	"regexp"
	"strconv"

	"github.com/gookit/color"
)

var hexRegexp = regexp.MustCompile("^#[0-9A-f]{6}$")

func GetPrinter(c string) color.PrinterFace {
	if IsHex(c) {
		return color.HEX(c)
	}

	cInt, err := strconv.Atoi(c)
	if err == nil {
		return color.S256(uint8(cInt))
	}

	if printer, ok := color.FgColors[c]; ok {
		return printer
	}

	if printer, ok := color.ExFgColors[c]; ok {
		return printer
	}

	return color.Normal
}

// Checks if the provided string is a valid hexadecimal color.
func IsHex(s string) bool {
	return hexRegexp.MatchString(s)
}

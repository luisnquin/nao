package utils_test

import (
	"testing"

	"github.com/luisnquin/nao/v3/internal/utils"
)

func TestSizeToStorageUnits(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{1500000, "1.43MB"},
		{2000000000, "1.86GB"},
		{5000, "4.88KB"},
		{0, "0.00KB"},
		{999, "0.98KB"},
		{100000000000, "93.13GB"},
		{-100, "-0.10KB"},
	}

	for _, test := range tests {
		result := utils.SizeToStorageUnits(test.input)
		if result != test.expected {
			t.Errorf("For input %d, expected %s, but got %s", test.input, test.expected, result)
		}
	}
}

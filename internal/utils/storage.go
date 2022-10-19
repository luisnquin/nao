package utils

import "fmt"

func BytesToStorageUnits(n int64) string {
	r := struct {
		parsedSize  float64
		storageUnit string
	}{}

	switch {
	case n > 1000000:
		r.parsedSize = float64(n) / (1 << 20)
		r.storageUnit = "MB"

	case n > 1_000_000_000:
		r.parsedSize = float64(n) / (1 << 30)
		r.storageUnit = "GB"

	default:
		r.parsedSize = float64(n) / (1 << 10)
		r.storageUnit = "KB"
	}

	return fmt.Sprintf("%.2f%s", r.parsedSize, r.storageUnit)
}

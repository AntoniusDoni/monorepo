package handler

import (
	"strconv"
)

// parsePositiveInt parses a string to a positive integer
func parsePositiveInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if val <= 0 {
		return 0, strconv.ErrRange
	}
	return val, nil
}

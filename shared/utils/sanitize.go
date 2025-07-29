package utils

import (
	"regexp"
	"strings"
)

// SanitizeSearchTerm removes special characters and normalizes spacing
func SanitizeSearchTerm(input string) string {
	re := regexp.MustCompile(`[^\w\s\-,.@]`)
	clean := re.ReplaceAllString(input, "")
	return strings.Join(strings.Fields(clean), " ")
}

// SanitizeAlphaNumeric keeps only letters and numbers
func SanitizeAlphaNumeric(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.ReplaceAllString(input, "")
}

// SanitizeUUID strips non-UUID characters
func SanitizeUUID(input string) string {
	re := regexp.MustCompile(`[^a-fA-F0-9\-]`)
	return re.ReplaceAllString(input, "")
}

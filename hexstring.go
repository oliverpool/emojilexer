package emojilexer

import (
	"strconv"
	"strings"
)

// FromHexString parses a dash-separated hexadecimal string to a string.
func FromHexString(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	var bb strings.Builder
	for _, r := range strings.Split(s, "-") {
		i, err := strconv.ParseInt(r, 16, 64)
		if err != nil {
			return bb.String(), err
		}
		bb.WriteRune(rune(i))
	}
	return bb.String(), nil
}

// ToHexString format a string as dash-separated hexadecimal string.
func ToHexString(s string) string {
	var bb strings.Builder
	for _, r := range s {
		if bb.Len() > 0 {
			bb.WriteRune('-')
		}
		bb.WriteString(strconv.FormatInt(int64(r), 16))
	}
	return bb.String()
}

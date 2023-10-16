package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func toByteSize(s string) (int64, error) {
	if s == "" {
		s = "64k"
	}
	s = strings.ToLower(s)
	var validSize = regexp.MustCompile(`[0-9]+[kmgtp]?`)
	if validSize.MatchString(s) {
		var mult int64 = 1
		if strings.HasSuffix(s, "k") {
			mult = 1024
			s = strings.TrimSuffix(s, "k")
		} else if strings.HasSuffix(s, "m") {
			mult = 1024 * 1024
			s = strings.TrimSuffix(s, "m")
		} else if strings.HasSuffix(s, "g") {
			mult = 1024 * 1024 * 1024
			s = strings.TrimSuffix(s, "g")
		} else if strings.HasSuffix(s, "t") {
			mult = 1024 * 1024 * 1024 * 1024
			s = strings.TrimSuffix(s, "t")
		} else if strings.HasSuffix(s, "p") {
			mult = 1024 * 1024 * 1024 * 1024 * 1024
			s = strings.TrimSuffix(s, "p")
		}
		v, err := strconv.Atoi(s)
		return int64(v) * mult, err
	}
	return 0, fmt.Errorf("invalid size value:%s", s)
}

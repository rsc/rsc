package main

import (
	"math"
	"strings"
)

func expandpkg(t0 string, pkg string) string {
	return strings.Replace(t0, `"".`, pkg+".", -1)
}

func double2ieee(ieee *uint64, f float64) {
	*ieee = math.Float64bits(f)
}

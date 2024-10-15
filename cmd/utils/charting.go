package utils

import (
	"fmt"
	"slices"
	"strings"
)

func MaxStringWidth(xs []string) int {
	m := 0
	for _, x := range xs {
		if strings.Count(x, "") > m {
			m = strings.Count(x, "")
		}
	}

	return m - 1
}

func TickScale(xs []float64, nTicks int) ([]string, int) {
	largest := int(slices.Max(xs))
	step := (largest / (nTicks - 1))

	var minStep int
	if step < 10 {
		minStep = 1
	} else if step < 50 {
		minStep = 5
	} else if step < 100 {
		minStep = 10
	} else {
		minStep = 100
	}

	for step*(nTicks-1) < largest || step%minStep != 0 {
		step++
	}

	scale := make([]string, nTicks)
	for i := 0; i < nTicks; i++ {
		scale[i] = fmt.Sprintf("%d", i*step)
	}

	return scale, step
}

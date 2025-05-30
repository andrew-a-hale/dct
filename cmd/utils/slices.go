package utils

import (
	"fmt"
	"maps"
	"math"
	"math/big"
	"slices"
)

type Vec2[K, V comparable] struct {
	X K
	Y V
}

func SortMap[K comparable](m map[K]int, dir int) []Vec2[K, int] {
	c := make(map[K]int)
	maps.Copy(c, m)

	var arr []Vec2[K, int]
	for len(c) > 0 {
		var max int
		var maxKey K
		for k, v := range c {
			if v > max {
				max = v
				maxKey = k
			}
		}

		arr = append(arr, Vec2[K, int]{maxKey, max})
		delete(c, maxKey)
	}

	return arr
}

type Summary struct {
	Min    int
	Mean   big.Float
	Median float64
	Max    int
	Count  int
	Sum    big.Int
}

func (s *Summary) String() string {
	return fmt.Sprint(s)
}

func Summarise(arr []int) Summary {
	var summary Summary
	summary.Min = math.MaxInt
	summary.Max = math.MinInt

	for _, v := range arr {
		summary.Count++
		summary.Sum.Add(&summary.Sum, big.NewInt(int64(v)))
		if v > summary.Max {
			summary.Max = v
		} else if v < summary.Min {
			summary.Min = v
		}
	}

	sum := new(big.Float).SetInt(&summary.Sum)
	div := big.NewFloat(float64(summary.Count))

	summary.Mean = *sum.Quo(sum, div)

	slices.Sort(arr)
	mid := len(arr) / 2
	if len(arr)%2 == 0 {
		summary.Median = float64(arr[mid]+arr[mid+1]) / 2
	} else {
		summary.Median = float64(arr[mid+1])
	}

	return summary
}

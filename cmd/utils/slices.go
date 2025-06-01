package utils

import (
	"fmt"
	"math"
	"sort"
	"unicode"
)

type Vec2[K, V comparable] struct {
	X K
	Y V
}

func SortMap[K comparable](m map[K]int, dir int) []Vec2[K, int] {
	var arr []Vec2[K, int]
	for k, v := range m {
		arr = append(arr, Vec2[K, int]{k, v})
	}

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Y > arr[j].Y
	})

	return arr
}

type Summary struct {
	Min   int
	Mean  float64
	Max   int
	Count int
	Sum   int
}

func (s Summary) String() string {
	return fmt.Sprintf("Min: %d\nMean: %f\nMax: %d", s.Min, s.Mean, s.Max)
}

func Summarise(m map[string]int) Summary {
	var summary Summary
	summary.Min = math.MaxInt
	summary.Max = math.MinInt

	for k := range m {
		length := len(k)
		summary.Count++
		summary.Sum += length
		if length > summary.Max {
			summary.Max = length
		}
		if length < summary.Min {
			summary.Min = length
		}
	}

	Assert(summary.Min < math.MaxInt, "error finding min string length")
	Assert(summary.Max > math.MinInt, "error finding max string length")

	summary.Mean = float64(summary.Sum) / float64(summary.Count)

	return summary
}

type Analysis struct {
	Control            int
	Comma              int
	Pipe               int
	Quotes             int
	Space              int
	NonSpaceWhitespace int
	NonAscii           int
	Rest               int
}

func (a Analysis) String() string {
	return fmt.Sprintf(`Control: %d
Comma: %d
Pipe: %d
Quotes: %d
Nonspace-Whitespace: %d
NonAscii: %d
Rest: %d`, a.Control, a.Comma, a.Pipe, a.Quotes, a.NonSpaceWhitespace, a.NonAscii, a.Rest)
}

func AnalyseRunes(m map[rune]int) Analysis {
	var analysis Analysis
	for k, v := range m {
		switch {
		case k == ' ':
			analysis.Space += v
		case unicode.IsSpace(k):
			analysis.NonSpaceWhitespace += v
		case unicode.IsControl(k):
			analysis.Control += v
		case k == ',':
			analysis.Comma += v
		case unicode.Is(unicode.Quotation_Mark, k):
			analysis.Quotes += v
		case k == '|':
			analysis.Pipe += v
		case k > unicode.MaxASCII:
			analysis.NonAscii += v
		default:
			analysis.Rest += v
		}
	}

	return analysis
}

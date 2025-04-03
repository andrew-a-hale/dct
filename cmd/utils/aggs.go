package utils

import (
	"log"
	"slices"
	"strings"
)

const (
	MEAN           = "mean"
	MEDIAN         = "median"
	MIN            = "min"
	MAX            = "max"
	SUM            = "sum"
	COUNT          = "count"
	COUNT_DISTINCT = "count_distinct"
)

var SUPPORTED_AGGS = []string{MEAN, MEDIAN, MIN, MAX, SUM, COUNT, COUNT_DISTINCT}

func CheckAgg(agg string) {
	if slices.Contains(SUPPORTED_AGGS, strings.ToLower(agg)) {
		return
	}

	log.Fatalf("agg '%s' not found, expect one of '%v'", agg, SUPPORTED_AGGS)
}

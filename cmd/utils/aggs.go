package utils

import (
	"log"
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
	for _, supportedAggs := range SUPPORTED_AGGS {
		if strings.ToLower(agg) == supportedAggs {
			return
		}
	}

	log.Fatalf("agg '%s' not found, expect one of '%v'", agg, SUPPORTED_AGGS)
}

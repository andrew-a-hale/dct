package utils

const (
	MEAN           = "mean"
	MEDIAN         = "median"
	MIN            = "min"
	MAX            = "max"
	COUNT_DISTINCT = "count_distinct"
)

var SUPPORTED_AGGS = []string{MEAN, MEDIAN, MIN, MAX, COUNT_DISTINCT}

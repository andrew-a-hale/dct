package utils

const (
	CSV     = ".csv"
	JSON    = ".json"
	NDJSON  = ".ndjson"
	PARQUET = ".parquet"
)

var PEEK_SUPPORTED_FILETYPES = []string{CSV, JSON, NDJSON, PARQUET}
var FLATTIFY_SUPPORTED_FILETYPES = []string{JSON, NDJSON}

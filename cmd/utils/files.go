package utils

const (
	CSV     = ".csv"
	JSON    = ".json"
	NDJSON  = ".ndjson"
	PARQUET = ".parquet"
	INVALID_FILE = "invalid"
)

var PEEK_SUPPORTED_FILETYPES = []string{CSV, JSON, NDJSON, PARQUET}
var FLATTIFY_SUPPORTED_FILETYPES = []string{JSON, NDJSON}
var PROFILE_SUPPORTED_FILETYPES = []string{CSV, JSON, NDJSON, PARQUET}

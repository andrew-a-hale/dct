package utils

import "fmt"

const (
	CSV          string = ".csv"
	JSON         string = ".json"
	NDJSON       string = ".ndjson"
	PARQUET      string = ".parquet"
	INVALID_FILE string = "invalid"
)

var (
	PEEK_SUPPORTED_FILETYPES     = []string{CSV, JSON, NDJSON, PARQUET}
	INFER_SUPPORTED_FILETYPES    = []string{CSV, JSON, NDJSON, PARQUET}
	PROFILE_SUPPORTED_FILETYPES  = []string{CSV, JSON, NDJSON, PARQUET}
	FLATTIFY_SUPPORTED_FILETYPES = []string{JSON, NDJSON}
)

type UnsupportedFileTypeErr struct {
	Msg      string
	Filename string
	Ext      string
}

func (e UnsupportedFileTypeErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Msg, e.Ext)
}

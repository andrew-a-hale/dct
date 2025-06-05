select
	json.values[0].a::decimal
	, json.values[0].b.a::decimal
	, json.values[1].a::decimal
	, json.values[1].b[0].values[0]::decimal
	, json.values[1].b[0].values[1]::decimal
	, json.values[1].b[0].values[2]::decimal
	, json.values[1].b[1].values[0]::decimal
	, json.values[1].b[1].values[1]::decimal
	, json.values[1].b[1].values[2]::decimal
	, json.values[2].a::decimal
	, json.values[2].b::decimal
	, json.values[3].a::decimal
	, json.values[3].b::decimal
	, json.values[4].a::decimal
	, json.values[4].b::decimal
from read_json_objects('test/resources/flattify.ndjson', format='unstructured') as json
